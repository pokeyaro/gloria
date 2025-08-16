package gloria

import (
	"fmt"
)

// Echo prints a one-line summary of the last request/response.
// Format: [API] METHOD URL -> STATUS (DURATION)
func (c *Client[T]) Echo() {
	fmt.Printf("[API] %s %s -> %d (%s)\n",
		c.meta.Method, c.meta.URL, c.meta.Status, c.meta.Duration)
}

// PostEchoHook returns an AfterHook that prints Echo() after SendCtx.
func PostEchoHook[T any]() AfterHook[T] {
	return func(c *Client[T]) error {
		c.Echo()
		return nil
	}
}
