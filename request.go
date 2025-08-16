package gloria

import (
	"net/url"
)

// SetHeader sets or overrides a single request header.
func (c *Client[T]) SetHeader(key, value string) *Client[T] {
	if c.hdr == nil {
		c.hdr = make(map[string]string)
	}
	c.hdr[key] = value
	return c
}

// SetHeaders sets or overrides multiple request headers.
func (c *Client[T]) SetHeaders(headers map[string]string) *Client[T] {
	if c.hdr == nil {
		c.hdr = make(map[string]string)
	}
	for k, v := range headers {
		c.hdr[k] = v
	}
	return c
}

// SetQueryParam sets or overrides a single query parameter.
func (c *Client[T]) SetQueryParam(key, value string) *Client[T] {
	if c.q == nil {
		c.q = make(map[string]string)
	}
	c.q[key] = value
	return c
}

// SetQueryParams sets or overrides multiple query parameters.
func (c *Client[T]) SetQueryParams(params map[string]string) *Client[T] {
	if c.q == nil {
		c.q = make(map[string]string)
	}
	for k, v := range params {
		c.q[k] = v
	}
	return c
}

// buildQuery encodes c.q into a URL-encoded query string without leading '?'.
func (c *Client[T]) buildQuery() string {
	if len(c.q) == 0 {
		return ""
	}
	v := url.Values{}
	for k, val := range c.q {
		v.Set(k, val)
	}
	return v.Encode()
}
