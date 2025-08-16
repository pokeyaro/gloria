package gloria

import (
	"encoding/base64"
)

// SetBasicAuth sets the Authorization header using HTTP Basic auth.
//
// Equivalent to: Authorization: Basic base64(username:password)
func (c *Client[T]) SetBasicAuth(username, password string) *Client[T] {
	if c.hdr == nil {
		c.hdr = &header{extra: make(map[string]string)}
	}
	cred := username + ":" + password
	c.hdr.extra["Authorization"] = "Basic " + base64.StdEncoding.EncodeToString([]byte(cred))
	return c
}

// SetBearerAuth sets the Authorization header using Bearer token.
//
// Equivalent to: Authorization: Bearer <token>
func (c *Client[T]) SetBearerAuth(token string) *Client[T] {
	if c.hdr == nil {
		c.hdr = &header{extra: make(map[string]string)}
	}
	c.hdr.extra["Authorization"] = "Bearer " + token
	return c
}
