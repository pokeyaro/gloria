package gloria

import (
	"time"
)

// Option mutates Config at construction time via NewClient(..., opts...).
type Option func(*Config)

// WithBaseURL sets the base URL used by SetRequest when path is relative.
func WithBaseURL(u string) Option {
	return func(c *Config) {
		if u != "" {
			c.BaseURL = u
		}
	}
}

// WithTimeout sets the default HTTP timeout.
func WithTimeout(d time.Duration) Option {
	return func(c *Config) {
		if d > 0 {
			c.Timeout = d
		}
	}
}

// WithREST toggles REST-style envelope decoding by default.
func WithREST(enabled bool) Option {
	return func(c *Config) { c.REST = enabled }
}

// WithOkCode sets the business success code for REST mode (default 0).
func WithOkCode(code int) Option {
	return func(c *Config) { c.OkCode = code }
}
