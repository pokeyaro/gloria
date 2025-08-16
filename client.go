package gloria

import (
	"net/http"
	"time"
)

// Client is the primary type for sending HTTP requests and decoding responses.
type Client[T any] struct {
	// config & codec
	cfg   Config
	codec Codec
	rest  RESTAdapter // nil => when REST is enabled, use DefaultRESTAdapter

	// builder state
	meta Meta
	hdr  *header           // request headers
	q    map[string]string // query parameters
	body []byte            // payload

	// hooks
	pre  []BeforeHook[T] // registered before-send hooks
	post []AfterHook[T]  // registered after-receive hooks

	// prepared request
	req *http.Request

	// transport
	httpc *http.Client // optional custom http.Client; if nil, SendCtx builds a default one

	// result state
	raw  []byte
	data T
	err  error
}

// Config defines settings for the client.
type Config struct {
	BaseURL string
	Timeout time.Duration

	// REST toggles REST-style envelope decoding.
	REST bool
	// OkCode is the REST success code (defaults to 0).
	OkCode int
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
			REST:    false,
			OkCode:  0,
		},
		codec: stdJSONCodec{}, // default codec
		rest:  nil,            // use default adapter when REST is enabled

		// builder state
		hdr: &header{
			// cookies: nil (append on nil is safe)
			extra: make(map[string]string),
		},
		q: make(map[string]string),
		// body: nil
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
