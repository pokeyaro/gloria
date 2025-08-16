package gloria

import (
	"net/http"
	"net/url"
)

// SetHeader sets or overrides a single request header.
func (c *Client[T]) SetHeader(key, value string) *Client[T] {
	if c.hdr == nil {
		c.hdr = &header{extra: make(map[string]string)}
	}
	if c.hdr.extra == nil {
		c.hdr.extra = make(map[string]string)
	}
	c.hdr.extra[key] = value
	return c
}

// SetHeaders sets or overrides multiple request headers.
func (c *Client[T]) SetHeaders(headers map[string]string) *Client[T] {
	if c.hdr == nil {
		c.hdr = &header{extra: make(map[string]string)}
	}
	if c.hdr.extra == nil {
		c.hdr.extra = make(map[string]string)
	}
	for k, v := range headers {
		c.hdr.extra[k] = v
	}
	return c
}

// SetAccept sets the "Accept" header value.
func (c *Client[T]) SetAccept(v string) *Client[T] {
	if c.hdr == nil {
		c.hdr = &header{extra: make(map[string]string)}
	}
	c.hdr.accept = v
	return c
}

// SetContentType sets the "Content-Type" header value.
func (c *Client[T]) SetContentType(v string) *Client[T] {
	if c.hdr == nil {
		c.hdr = &header{extra: make(map[string]string)}
	}
	c.hdr.contentType = v
	return c
}

// SetLanguage sets the "Accept-Language" header value.
func (c *Client[T]) SetLanguage(v string) *Client[T] {
	if c.hdr == nil {
		c.hdr = &header{extra: make(map[string]string)}
	}
	c.hdr.language = v
	return c
}

// SetUserAgent sets the "User-Agent" header value.
func (c *Client[T]) SetUserAgent(v string) *Client[T] {
	if c.hdr == nil {
		c.hdr = &header{extra: make(map[string]string)}
	}
	c.hdr.userAgent = v
	return c
}

// SetCookie adds a cookie to the request.
func (c *Client[T]) SetCookie(ck *http.Cookie) *Client[T] {
	if c.hdr == nil {
		c.hdr = &header{extra: make(map[string]string)}
	}
	c.hdr.cookies = append(c.hdr.cookies, ck)
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
