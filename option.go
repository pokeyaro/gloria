package gloria

import (
	"time"
)

// Option applies a configuration change to Config.
type Option func(*Config)

// WithBaseURL sets the base URL for the client.
func WithBaseURL(baseURL string) Option {
	return func(cfg *Config) {
		cfg.BaseURL = baseURL
	}
}

// WithTimeout sets the request timeout for the client.
func WithTimeout(d time.Duration) Option {
	return func(cfg *Config) {
		cfg.Timeout = d
	}
}
