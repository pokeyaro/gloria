package transport

import (
	"net/http"
	"time"
)

// Config carries minimal knobs for building the default HTTP client.
// (Add TLS/proxy/retry knobs here later.)
type Config struct {
	Timeout time.Duration
}

// Build returns a default *http.Client for Gloria.
// Keep it tiny for now; we can swap Transport, TLS, proxy, retry later.
func Build(cfg Config) *http.Client {
	return &http.Client{
		Timeout: cfg.Timeout,
	}
}
