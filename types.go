package gloria

import (
	"net/http"
)

// header groups request header fields for a request.
type header struct {
	accept      string
	contentType string
	language    string
	userAgent   string
	cookies     []*http.Cookie
	extra       map[string]string
}
