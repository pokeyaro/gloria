package gloria

import (
	"errors"
	"fmt"
)

// RESTAdapter unwraps a REST-style envelope and optionally decodes its data into out.
// It returns business code, message, and an error if decoding the envelope fails.
// Whether a code is considered OK is decided by IsOK.
type RESTAdapter interface {
	Unwrap(raw []byte, codec Codec, out any) (code int, msg string, err error)
	IsOK(code int) bool
}

// RESTError represents a REST-layer business failure returned by DecodeREST.
type RESTError struct {
	Code   int    // business code
	Msg    string // business message
	Status int    // HTTP status code (from last response)
	Raw    []byte // raw response body (for diagnostics)
	Cause  error  // underlying cause, if any
}

func (e *RESTError) Error() string {
	if e.Status > 0 {
		return fmt.Sprintf("rest error: code=%d msg=%s (http=%d)", e.Code, e.Msg, e.Status)
	}
	return fmt.Sprintf("rest error: code=%d msg=%s", e.Code, e.Msg)
}

// Unwrap returns the underlying cause to support errors.Is/As/Unwrap.
func (e *RESTError) Unwrap() error {
	return e.Cause
}

// IsRESTError reports whether err is a RESTError.
func IsRESTError(err error) bool {
	var e *RESTError
	return errors.As(err, &e)
}

// AsRESTError extracts *RESTError if present.
func AsRESTError(err error) (*RESTError, bool) {
	var e *RESTError
	if errors.As(err, &e) {
		return e, true
	}
	return nil, false
}

// DefaultRESTAdapter unwraps {"code": int, "msg": string, "data": any}.
// OkCode decides success (default 0), matching v1 behavior.
type DefaultRESTAdapter struct {
	OkCode int
}

func (a DefaultRESTAdapter) IsOK(code int) bool {
	return code == a.OkCode
}

// genericEnvelope is used internally by DefaultRESTAdapter.
type genericEnvelope struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

// Unwrap implements RESTAdapter using v1-compatible keys: code/msg/data.
func (a DefaultRESTAdapter) Unwrap(raw []byte, codec Codec, out any) (int, string, error) {
	var env genericEnvelope
	if err := codec.Unmarshal(raw, &env); err != nil {
		return 0, "", err
	}
	if out != nil {
		bs, err := codec.Marshal(env.Data)
		if err != nil {
			return env.Code, env.Msg, err
		}
		if err := codec.Unmarshal(bs, out); err != nil {
			return env.Code, env.Msg, err
		}
	}
	return env.Code, env.Msg, nil
}

// EnableREST enables REST-style envelope decoding.
func (c *Client[T]) EnableREST() *Client[T] {
	c.cfg.REST = true
	return c
}

// DisableREST disables REST-style envelope decoding.
func (c *Client[T]) DisableREST() *Client[T] {
	c.cfg.REST = false
	return c
}

// ToggleMode switches between HTTP mode (false) and REST mode (true).
func (c *Client[T]) ToggleMode() *Client[T] {
	c.cfg.REST = !c.cfg.REST
	return c
}

// DefineOkCode sets the REST success code (default 0).
func (c *Client[T]) DefineOkCode(code int) *Client[T] {
	c.cfg.OkCode = code
	return c
}

// UseRESTAdapter sets a custom RESTAdapter implementation and enables REST mode.
func (c *Client[T]) UseRESTAdapter(a RESTAdapter) *Client[T] {
	c.rest = a
	c.cfg.REST = true
	return c
}

// DecodeREST unwraps the REST envelope and validates business code.
// On success, it stores the decoded data into c.data (type T).
func (c *Client[T]) DecodeREST() (*Client[T], error) {
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

	adapter := c.rest
	if adapter == nil {
		adapter = DefaultRESTAdapter{OkCode: c.cfg.OkCode}
	}

	var dst T
	code, msg, err := adapter.Unwrap(c.raw, c.codec, &dst)
	if err != nil {
		c.err = err
		return c, err
	}
	if !adapter.IsOK(code) {
		c.err = &RESTError{
			Code:   code,
			Msg:    msg,
			Status: c.meta.Status,
			Raw:    c.raw,
		}
		return c, c.err
	}
	c.data = dst
	return c, nil
}

// DecodeAuto chooses REST or plain JSON decoding based on config.
func (c *Client[T]) DecodeAuto() (*Client[T], error) {
	if c.cfg.REST {
		return c.DecodeREST()
	}
	return c.Decode()
}
