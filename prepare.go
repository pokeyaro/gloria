package gloria

import (
	"bytes"
	"context"
	"io"
	"net/http"
)

// Prepare builds an *http.Request from the current client state without sending it.
// It applies method, URL (with query parameters), headers, cookies, and body.
func (c *Client[T]) Prepare(ctx context.Context) (*http.Request, error) {
	if c.err != nil {
		return nil, c.err
	}

	finalURL := c.meta.URL
	if q := c.buildQuery(); q != "" {
		if hasQuery(finalURL) {
			finalURL = finalURL + "&" + q
		} else {
			finalURL = finalURL + "?" + q
		}
	}

	// Use io.Reader interface to avoid passing a typed-nil concrete reader.
	var body io.Reader
	if len(c.body) > 0 {
		body = bytes.NewReader(c.body)
	} else {
		body = nil
	}

	req, err := http.NewRequestWithContext(ctx, c.meta.Method, finalURL, body)
	if err != nil {
		return nil, err
	}

	// Apply headers set via semantic setters and extra kv headers.
	if c.hdr != nil {
		if c.hdr.accept != "" {
			req.Header.Set("Accept", c.hdr.accept)
		}
		if c.hdr.contentType != "" {
			req.Header.Set("Content-Type", c.hdr.contentType)
		}
		if c.hdr.language != "" {
			req.Header.Set("Accept-Language", c.hdr.language)
		}
		if c.hdr.userAgent != "" {
			req.Header.Set("User-Agent", c.hdr.userAgent)
		}
		for k, v := range c.hdr.extra {
			req.Header.Set(k, v)
		}
		for _, ck := range c.hdr.cookies {
			if ck != nil {
				req.AddCookie(ck)
			}
		}
	}

	c.req = req
	return req, nil
}

// hasQuery reports whether s already contains '?'.
func hasQuery(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == '?' {
			return true
		}
	}
	return false
}
