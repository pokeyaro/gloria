package gloria

import (
	"context"

	ex "github.com/pokeyaro/gloria/v2/internal/exec"
	tr "github.com/pokeyaro/gloria/v2/internal/transport"
)

// SendCtx sends the prepared request (or prepares one if absent) and
// stores response metadata and raw body on the client.
func (c *Client[T]) SendCtx(ctx context.Context) (*Client[T], error) {
	if c.err != nil {
		return c, c.err
	}

	// pre hooks
	for _, h := range c.pre {
		if h == nil {
			continue
		}
		if err := h(c); err != nil {
			c.err = err
			return c, err
		}
	}

	// ensure request exists
	if c.req == nil {
		if _, err := c.Prepare(ctx); err != nil {
			c.err = err
			return c, err
		}
	}

	// choose client: prefer injected, else internal builder
	httpClient := c.httpc
	if httpClient == nil {
		httpClient = tr.Build(tr.Config{
			Timeout: c.cfg.Timeout,
		})
	}

	// execute (internal exec)
	res, err := ex.Do(ctx, httpClient, c.req)
	if err != nil {
		c.err = err
		return c, err
	}

	// record meta
	c.meta.Status = res.Status
	c.meta.Proto = res.Proto
	c.meta.Duration = res.Duration
	c.meta.ReceivedAt = res.ReceivedAt

	// keep raw body
	c.raw = res.Body

	// post hooks
	for _, h := range c.post {
		if h == nil {
			continue
		}
		if err := h(c); err != nil {
			c.err = err
			return c, err
		}
	}

	// Note: we do not treat non-2xx as an error here; the caller can decide.
	return c, nil
}

// Send is a convenience wrapper over SendCtx with context.Background().
func (c *Client[T]) Send() (*Client[T], error) {
	return c.SendCtx(context.Background())
}
