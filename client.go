package gloria

import (
	"time"
)

// Client is the primary type for sending HTTP requests and decoding responses.
type Client[T any] struct {
	cfg  Config
	meta Meta
	data T
	err  error
}

// Config defines settings for the client.
type Config struct {
	BaseURL string
	Timeout time.Duration
}

// Meta contains metadata about a request execution.
type Meta struct {
	Method   string
	URL      string
	Duration time.Duration
}

// NewClient creates a new Client with default configuration.
func NewClient[T any](baseURL string, opts ...Option) *Client[T] {
	c := &Client[T]{
		cfg: Config{
			BaseURL: baseURL,
			Timeout: 30 * time.Second,
		},
	}
	for _, opt := range opts {
		opt(&c.cfg)
	}
	return c
}

// New is an alias of NewClient.
func New[T any](baseURL string, opts ...Option) *Client[T] {
	return NewClient[T](baseURL, opts...)
}
