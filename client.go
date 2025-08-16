package gloria

import (
	"net/http"
	"time"
)

// Client is the primary type for sending HTTP requests and decoding responses.
type Client[T any] struct {
	cfg   Config
	codec Codec

	meta Meta
	hdr  *header           // request headers
	q    map[string]string // query parameters
	body []byte            // payload

	pre  []BeforeHook[T] // registered before-send hooks
	post []AfterHook[T]  // registered after-receive hooks

	req *http.Request

	raw  []byte
	data T
	err  error
}

// Config defines settings for the client.
type Config struct {
	BaseURL string
	Timeout time.Duration
}

// Meta records basic request/response metadata.
type Meta struct {
	Method     string
	URL        string
	Status     int
	Proto      string
	Duration   time.Duration
	ReceivedAt time.Time
}

// NewClient creates a new Client with default configuration.
func NewClient[T any](baseURL string, opts ...Option) *Client[T] {
	c := &Client[T]{
		cfg: Config{
			BaseURL: baseURL,
			Timeout: 30 * time.Second,
		},
		codec: stdJSONCodec{},
		hdr: &header{
			extra: make(map[string]string),
		},
		q: make(map[string]string),
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
