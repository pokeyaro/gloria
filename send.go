package gloria

import (
	"context"
	"io"
	"net/http"
	"time"
)

// SendCtx sends the prepared request (or prepares one if absent) and
// stores response metadata and raw body on the client.
func (c *Client[T]) SendCtx(ctx context.Context) (*Client[T], error) {
	if c.err != nil {
		return c, c.err
	}

	// Ensure we have a request.
	if c.req == nil {
		if _, err := c.Prepare(ctx); err != nil {
			c.err = err
			return c, err
		}
	}

	httpClient := &http.Client{
		Timeout: c.cfg.Timeout,
	}

	start := time.Now()
	resp, err := httpClient.Do(c.req)
	if err != nil {
		c.err = err
		return c, err
	}
	defer resp.Body.Close()

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		c.err = err
		return c, err
	}

	// Record response metadata.
	c.meta.Status = resp.StatusCode
	c.meta.Proto = resp.Proto
	c.meta.Duration = time.Since(start)
	c.meta.ReceivedAt = time.Now()

	// Preserve raw response body for later decoding.
	c.raw = bs

	// Note: we do not treat non-2xx as an error here; the caller can decide.
	return c, nil
}
