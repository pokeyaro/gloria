package gloria

// SetJSON marshals v with the current codec and sets it as the request body.
// If Content-Type is not set, it is set to "application/json".
func (c *Client[T]) SetJSON(v any) *Client[T] {
	if c.codec == nil {
		c.codec = stdJSONCodec{}
	}
	bs, err := c.codec.Marshal(v)
	if err != nil {
		c.err = err
		return c
	}
	c.body = bs
	if c.hdr == nil {
		c.hdr = &header{extra: make(map[string]string)}
	}
	if c.hdr.contentType == "" {
		c.hdr.contentType = "application/json"
	}
	return c
}

// SetBody sets raw bytes as the request body.
// The caller is responsible for setting an appropriate Content-Type.
func (c *Client[T]) SetBody(b []byte) *Client[T] {
	c.body = b
	return c
}
