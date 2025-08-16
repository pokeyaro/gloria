package gloria

// Raw returns the last response body bytes.
func (c *Client[T]) Raw() []byte {
	return c.raw
}

// Status returns the last HTTP status code.
func (c *Client[T]) Status() int {
	return c.meta.Status
}
