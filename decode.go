package gloria

import (
	"errors"
)

// Decode unmarshals the raw response body into the client's generic data (T)
// using the configured codec. It stores the result in c.data.
func (c *Client[T]) Decode() (*Client[T], error) {
	if c.err != nil {
		return c, c.err
	}
	if len(c.raw) == 0 {
		c.err = errors.New("empty response body")
		return c, c.err
	}
	if c.codec == nil {
		c.codec = stdJSONCodec{}
	}

	var dst T
	if err := c.codec.Unmarshal(c.raw, &dst); err != nil {
		c.err = err
		return c, err
	}
	c.data = dst
	return c, nil
}

// DecodeInto unmarshals the raw response body into v using the configured codec.
// v must be a non-nil pointer to a decodable value.
func (c *Client[T]) DecodeInto(v any) error {
	if c.err != nil {
		return c.err
	}
	if len(c.raw) == 0 {
		c.err = errors.New("empty response body")
		return c.err
	}
	if c.codec == nil {
		c.codec = stdJSONCodec{}
	}
	return c.codec.Unmarshal(c.raw, v)
}
