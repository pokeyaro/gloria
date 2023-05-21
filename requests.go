// Copyright (c) 2023 Pokeya Boa <pokeya.mystic@gmail.com>, All rights reserved.
// resty source code and usage is governed by a MIT style
// license that can be found in the LICENSE file.

package gloria

import (
	"net/http"
)

/*
	The following is the request method
*/

// request is a generic function used for sending HTTP requests with the specified
// method, path, parameters, data, and headers.
// It returns a new client instance configured for the request.
//
// The function performs the following steps:
// 1. Validates the method to ensure it is a valid HTTP method.
// 2. Parses the URL segments from the path.
// 3. Initializes a new client instance using the default settings.
// 4. Sets the request method for the client.
// 5. Sets the URL for the client based on the parsed URL segments.
// 6. Sets the query parameters for the client, unless the method is OPTIONS.
// 7. Sets the request payload (body) for the client, unless the method is GET or OPTIONS.
// 8. Sets the request headers for the client.
// 9. Sends the request using the client.
//
// The function returns a client instance configured for the request.
func request[T any](method, path string, params H, data any, headers ...H) *Client[T] {
	// Check if the method is valid
	isValidMethod(method)

	// Parse the URL
	parseUrl := urlSegments(path)

	// Initialize a new client
	r := Default[T]()

	// Set the request method
	r.SetMethod(method)

	// Set the URL
	r.SetURL(parseUrl.scheme, parseUrl.host, parseUrl.baseURI, parseUrl.endpoint)

	// Set the query parameters
	if method != MethodOptions {
		if isEmpty(parseUrl.params) {
			r.SetQueryParams(params)
		} else {
			r.SetQueryParams(parseUrl.params)
		}
	}

	// Set the request payload
	if method != MethodGet && method != MethodOptions {
		if isEmpty(data) {
			r.SetPayload(nil)
		} else {
			r.SetPayload(data)
		}
	}

	// Set the request headers
	if len(headers) > 0 {
		r.SetHeaders(headers[0])
	}

	// Send the request
	r.Send()

	return r
}

/*
	A Classic Request Style like python requests library
*/

// GET is a shorthand function for creating a GET request with the specified path,
// query parameters, and headers.
// It returns a new client instance configured for a GET request.
func GET[T any](path string, params H, headers ...H) *Client[T] {
	return request[T](http.MethodGet, path, params, Placeholder, headers...)
}

// POST is a shorthand function for creating a POST request with the specified path,
// query parameters, request body data, and headers.
// It returns a new client instance configured for a POST request.
func POST[T any](path string, params H, data any, headers ...H) *Client[T] {
	return request[T](http.MethodPost, path, params, data, headers...)
}

// PUT is a shorthand function for creating a PUT request with the specified path,
// query parameters, request body data, and headers.
// It returns a new client instance configured for a PUT request.
func PUT[T any](path string, params H, data any, headers ...H) *Client[T] {
	return request[T](http.MethodPut, path, params, data, headers...)
}

// DELETE is a shorthand function for creating a DELETE request with the specified
// path, query parameters, request body data, and headers.
// It returns a new client instance configured for a DELETE request.
func DELETE[T any](path string, params H, data any, headers ...H) *Client[T] {
	return request[T](http.MethodDelete, path, params, data, headers...)
}

// PATCH is a shorthand function for creating a PATCH request with the specified
// path, query parameters, request body data, and headers.
// It returns a new client instance configured for a PATCH request.
func PATCH[T any](path string, params H, data any, headers ...H) *Client[T] {
	return request[T](http.MethodPatch, path, params, data, headers...)
}

// HEAD is a shorthand function for creating a HEAD request with the specified
// path, query parameters, request body data, and headers.
// It returns a new client instance configured for a HEAD request.
func HEAD[T any](path string, params H, headers ...H) *Client[T] {
	return request[T](http.MethodHead, path, params, Placeholder, headers...)
}

// OPTIONS is a shorthand function for creating an OPTIONS request with the specified
// path and headers.
// It returns a new client instance configured for an OPTIONS request.
func OPTIONS[T any](path string, headers ...H) *Client[T] {
	return request[T](http.MethodOptions, path, nil, Placeholder, headers...)
}

/*
	A Different Request Style Choice
*/

type ExecMethod[T any] func(method string) *RESTFulResp[T]

// Request is a helper function that creates an ExecMethod function for making HTTP requests.
// It takes the path, query parameters, request payload, and optional headers as input.
// It returns an ExecMethod function that can be used to execute the specified HTTP method.
func Request[T any](path string, params H, data any, headers ...H) ExecMethod[T] {
	return func(method string) *RESTFulResp[T] {
		// Return restful response
		return request[T](method, path, params, data, headers...).Result
	}
}
