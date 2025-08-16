package gloria

// BeforeHook is executed before the HTTP request is sent.
// It may mutate the client; returning a non-nil error aborts the send.
type BeforeHook[T any] func(*Client[T]) error

// AfterHook is executed after the HTTP response is received and the body is read.
// It may mutate the client; returning a non-nil error marks the send as failed.
type AfterHook[T any] func(*Client[T]) error

// UsePreHooks registers one or more before-send hooks.
func (c *Client[T]) UsePreHooks(funcs ...BeforeHook[T]) *Client[T] {
	c.pre = append(c.pre, funcs...)
	return c
}

// UsePostHooks registers one or more after-receive hooks.
func (c *Client[T]) UsePostHooks(funcs ...AfterHook[T]) *Client[T] {
	c.post = append(c.post, funcs...)
	return c
}
