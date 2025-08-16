package gloria

import (
	"strings"
)

// Common HTTP method constants.
const (
	MethodGet     = "GET"
	MethodPost    = "POST"
	MethodPut     = "PUT"
	MethodDelete  = "DELETE"
	MethodPatch   = "PATCH"
	MethodHead    = "HEAD"
	MethodOptions = "OPTIONS"
)

// SetRequest sets the HTTP method and request URL on the client.
// If url is relative and BaseURL is set, the final URL is BaseURL + url.
// Method is normalized to upper-case and must be one of the standard constants.
func (c *Client[T]) SetRequest(method, url string) *Client[T] {
	m := strings.ToUpper(strings.TrimSpace(method))
	if !isStdMethod(m) {
		panic("gloria: unsupported HTTP method: " + m)
	}
	c.meta.Method = m
	c.meta.URL = resolveURL(c.cfg.BaseURL, url)
	return c
}

// isStdMethod reports whether m is a supported HTTP verb.
func isStdMethod(m string) bool {
	switch m {
	case MethodGet, MethodPost, MethodPut, MethodDelete,
		MethodPatch, MethodHead, MethodOptions:
		return true
	default:
		return false
	}
}

// resolveURL joins base and path if path is relative.
// If path is absolute (starts with "http://"/"https://"), it is returned as-is.
func resolveURL(base, path string) string {
	p := strings.TrimSpace(path)
	if hasScheme(p) || base == "" {
		return p
	}
	b := strings.TrimRight(base, "/")
	pp := "/" + strings.TrimLeft(p, "/")
	return b + pp
}

// hasScheme reports whether s starts with an HTTP(S) scheme prefix.
func hasScheme(s string) bool {
	ss := strings.ToLower(s)
	return strings.HasPrefix(ss, "http://") || strings.HasPrefix(ss, "https://")
}
