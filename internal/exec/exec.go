package exec

import (
	"context"
	"io"
	"net/http"
	"time"
)

type Result struct {
	Status     int
	Proto      string
	Body       []byte
	Duration   time.Duration
	ReceivedAt time.Time
}

// Do performs the HTTP round trip and reads the entire response body.
func Do(ctx context.Context, hc *http.Client, req *http.Request) (*Result, error) {
	start := time.Now()

	// ensure request carries ctx (in case caller built it without)
	req = req.WithContext(ctx)

	resp, err := hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &Result{
		Status:     resp.StatusCode,
		Proto:      resp.Proto,
		Body:       bs,
		Duration:   time.Since(start),
		ReceivedAt: time.Now(),
	}, nil
}
