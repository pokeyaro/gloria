package gloria

// SetBody sets the raw request payload bytes.
// It copies the input to avoid retaining external slices.
func (c *Client[T]) SetBody(b []byte) *Client[T] {
	if b == nil {
		c.body = nil
		return c
	}
	c.body = append(c.body[:0], b...)
	return c
}

// SetJSON marshals v with the current Codec and sets it as the request body.
// It also sets "Content-Type: application/json" if not already set.
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

	// ensure Content-Type
	if c.hdr == nil {
		c.hdr = &header{extra: make(map[string]string)}
	}
	if c.hdr.contentType == "" {
		c.hdr.contentType = "application/json"
	}
	return c
}

// ResetBody clears the prepared request payload.
func (c *Client[T]) ResetBody() *Client[T] {
	c.body = nil
	return c
}
