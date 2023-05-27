// Copyright (c) 2023 Pokeya Boa <pokeya.mystic@gmail.com>, All rights reserved.
// resty source code and usage is governed by a MIT style
// license that can be found in the LICENSE file.

package gloria

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

/*
Constructor to create a client instance
*/

// New function returns an empty template initialization (for rest mode).
func New[T any]() *Client[T] {
	client := &Client[T]{
		Context: &Context{
			HttpClient: &http.Client{},
			Request:    &http.Request{},
			Response:   &Response{},
		},
		Meta: &Meta{},
		Config: &Config{
			FilterSlash:   false,
			IsDebug:       false,
			Logger:        nil,
			IsRestMode:    true,
			DefaultOkCode: OkCode,
			JSONLoader:    NativeJSONLibrary{},
		},
		Exception:     &Exception{},
		Result:        &RESTFulResp[T]{},
		beforeRequest: []func(*Client[T]) error{},
		afterResponse: []func(*Client[T]) error{},
		urls:          &urls{},
		params:        SMap{},
		authorization: &authorization{},
		headers: &header{
			cookies: []*http.Cookie{},
			extra:   SMap{},
		},
		payload: nil,
	}

	return client
}

// Default function returns a basic default template.
//
// Preset parameters:
//  1. Debug: false
//  2. Logger: true
//  3. OkCode: 0
//  4. TLS: skip
//  5. Timeout: 30s
//  6. RestMode: true
//
// Inject default middleware:
//  1. set Accept: application/json
//  2. set Content-Type: application/json
//  3. set Content-Language: en-US,en;q=0.9
//  4. set User-Agent: - actual environment
func Default[T any]() *Client[T] {
	// Create an empty client object
	client := New[T]()

	// Add default option parameter
	client.Optional(
		// single setting
		WithIsDebug[T](false),
		WithUseLogger[T](true),
		Lambda[T](func(c *Client[T]) {
			// More settings
			c.Config.SkipTLS = true
			c.Config.Timeout = TimeoutMedium
			c.Config.IsRestMode = true
			c.Config.DefaultOkCode = OkCode
			c.Config.JSONLoader = GoJSONLibrary{}
		}),
	)

	// Add hook action (load default request middleware)
	client.UsePreHooks(func(c *Client[T]) error {
		if isEmpty(c.headers) {
			c.headers = &header{
				accept:      JsonContentType,
				contentType: JsonContentType,
				language:    LocaleEn,
				userAgent:   getUserAgent(),
			}
		}

		return nil
	})

	return client
}

// NewREST function returns an empty template initialization (for rest mode).
func NewREST[T any]() *Client[T] {
	client := New[T]()

	return client
}

// NewHTTP function returns an empty template initialization (for http mode).
func NewHTTP[T any]() *Client[T] {
	client := New[T]().ToggleMode()

	return client
}

// Deprecated: NewClient function creates an empty template similar to the New function.
// However, it is not recommended as it no longer has generic type features.
func NewClient() *Client[any] {
	client := New[any]()

	return client
}

/*
	Exposed function optional methods for the Client struct
*/

type ClientFunc[T any] func(*Client[T])

// Optional is a method that allows applying optional configurations to the client instance.
// It accepts a variadic parameter fns, which represents a list of client functions to be applied.
// Each function takes the client instance as a parameter and applies specific configurations or
// modifications to it.
// The method iterates over the list of functions and calls each function with the client instance
// as an argument.
// After applying all the functions, it returns the modified client instance.
func (c *Client[T]) Optional(fns ...ClientFunc[T]) *Client[T] {
	for _, fn := range fns {
		fn(c)
	}
	return c
}

// Lambda is a helper function that converts a function f into a ClientFunc[T].
// It takes a function f as a parameter, which accepts a client instance and performs specific
// actions on it.
// The Lambda function wraps the provided function f and returns it as a ClientFunc[T].
// The returned ClientFunc[T] can be used as an argument in the Optional method to apply the
// provided function to the client instance.
func Lambda[T any](f func(*Client[T])) ClientFunc[T] {
	return f
}

// WithTimeout is a ClientFunc[T] function that configures the timeout duration for a client
// instance.
// It takes a time.Duration value timeout as a parameter and returns a ClientFunc[T].
// When applied to a client instance using the Optional method, it configures the timeout
// duration based on the provided value.
// The timeout duration is capped to a minimum of TimeoutShort and a maximum of TimeoutLong.
func WithTimeout[T any](timeout time.Duration) ClientFunc[T] {
	return func(c *Client[T]) {
		if timeout < TimeoutShort {
			timeout = TimeoutShort
		}
		if timeout > TimeoutLong {
			timeout = TimeoutLong
		}

		c.Config.Timeout = timeout
	}
}

// WithSkipTLS is a ClientFunc[T] function that sets the SkipTLS configuration of a client
// instance.
// It takes a boolean value skipTLS as a parameter and returns a ClientFunc[T].
// When applied to a client instance using the Optional method, it sets the SkipTLS
// configuration of the client to the provided value.
func WithSkipTLS[T any](skipTLS bool) ClientFunc[T] {
	return func(c *Client[T]) {
		c.Config.SkipTLS = skipTLS
	}
}

// WithIsDebug is a ClientFunc[T] function that sets the IsDebug configuration of a client instance.
// It takes a boolean value isDebug.
func WithIsDebug[T any](isDebug bool) ClientFunc[T] {
	return func(c *Client[T]) {
		c.Config.IsDebug = isDebug
	}
}

// WithUseLogger is a ClientFunc[T] function that configures the usage of a logger for a
// client instance.
// It takes a boolean value enabled as a parameter and returns a ClientFunc[T].
// When applied to a client instance using the Optional method, it configures the usage
// of a logger based on the provided value.
// If enabled is true, it configures a custom log formatter for the client's logger.
// If enabled is false, it does not configure a logger for the client.
func WithUseLogger[T any](enabled bool) ClientFunc[T] {
	return func(c *Client[T]) {
		if enabled {
			// Configure a custom log formatter for the Logger.
			// logger := log.New(os.Stdout, "", log.Lshortfile|log.Ldate|log.Ltime)
			logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
			c.Config.Logger = logger
		}
	}
}

// WithRegisterJsonLibrary is a ClientFunc[T] function that registers the json library for a
// client instance.
// You can choose the popular json parsing library independently.
// Note: Please implement the JSONLibrary interface definition yourself.
func WithRegisterJsonLibrary[T any](lib JSONLibrary) ClientFunc[T] {
	return func(c *Client[T]) {
		c.Config.JSONLoader = lib
	}
}

// Deprecated: WithFilterSlash is a ClientFunc[T] function that sets the FilterSlash configuration of a client instance.
// It takes a boolean parameter filterSlash to enable or disable filtering of trailing slashes in URLs.
// When filterSlash is set to true, the client will remove any trailing slashes from the URLs it sends requests to.
// This can be useful in cases where the server treats URLs with and without trailing slashes differently.
// Note: that filtering slashes in URLs may affect the behavior of your requests, so use it carefully.
func WithFilterSlash[T any](filterSlash bool) ClientFunc[T] {
	return func(c *Client[T]) {
		c.Config.FilterSlash = filterSlash
	}
}

// Deprecated: WithDisabledRestMode Close the restful api mode, the request body will become the content of the T generic.
// Notes: This method is useful if your response body is not a standard restful response format data!
// Instead: Please use the ToggleMode method.
func WithDisabledRestMode[T any]() ClientFunc[T] {
	return func(c *Client[T]) {
		c.Config.IsRestMode = false
	}
}

// Deprecated: WithModifySuccessCode sets the success return code for the client. The success code is used to identify
// successful responses. If not explicitly set, the default success code is 0.
func WithModifySuccessCode[T any](code int) ClientFunc[T] {
	return func(c *Client[T]) {
		c.Config.DefaultOkCode = code
	}
}

/*
	A shortcut
*/

func (c *Client[T]) ToggleMode() *Client[T] {
	c.Config.IsRestMode = !c.Config.IsRestMode

	return c
}

// Note: The FilterUrlSlash method needs to be called before the SetURL or SetRequest method.
func (c *Client[T]) FilterUrlSlash() *Client[T] {
	c.Config.FilterSlash = true

	return c
}

func (c *Client[T]) DefineOkCode(code int) *Client[T] {
	c.Config.DefaultOkCode = code

	return c
}

// Note: Please implement the JSONLibrary interface definition yourself.
func (c *Client[T]) RegisterJsonLib(lib JSONLibrary) *Client[T] {
	c.Config.JSONLoader = lib

	return c
}

/*
	Exposed chain methods with Getter attribute for the Client struct
*/

// GetQuery returns the value of the specified query parameter from the client instance.
func (c *Client[T]) GetQuery(q string) string {
	qs := c.params
	if isEmpty(qs) {
		return ""
	}

	return qs[q]
}

// GetQueryParams returns the query parameters as a SMap from the client instance.
func (c *Client[T]) GetQueryParams() SMap {
	qs := c.params
	if isEmpty(qs) {
		return nil
	}

	return qs
}

// GetHeader returns the value of the specified header key from the client's request context.
func (c *Client[T]) GetHeader(key string) string {
	hdr := c.Context.Request.Header.Get(key)

	return hdr
}

// GetHeaders returns the headers as a http.Header from the client's request context.
func (c *Client[T]) GetHeaders() http.Header {
	headers := c.Context.Request.Header.Clone()
	if isEmpty(headers) {
		return nil
	}

	return headers
}

// GetCookie returns the cookie with the specified name from the client's request context.
// If the cookie is found, it returns the cookie object.
// If the cookie is not found, it returns an error.
func (c *Client[T]) GetCookie(name string) (*http.Cookie, error) {
	return c.Context.Request.Cookie(name)
}

// GetCookies returns the cookies as a []*http.Cookie from the client's request context.
func (c *Client[T]) GetCookies() []*http.Cookie {
	cookies := c.Context.Request.Cookies()
	if isEmpty(cookies) {
		return nil
	}

	return cookies
}

/*
	Exposed chain methods with Setter attribute for the Client struct
*/

// SetMethod sets the HTTP method for the client instance to the specified method.
// This method must be one of the supported methods: GET, POST, PUT, PATCH, DELETE, HEAD, or OPTIONS.
//
// If an unsupported method is provided, this method will panic with an error message indicating the supported methods.
//
// Example usage:
//
//	client.SetMethod("GET")
func (c *Client[T]) SetMethod(method string) *Client[T] {
	method = strings.ToUpper(method)

	isValidMethod(method)

	c.Meta.Method = method

	return c
}

// SetURL sets the scheme, host, base URI, and endpoint for the client instance.
// These values will be used to construct the URL for the request.
//
// The scheme parameter specifies the protocol scheme for the request, such as "http" or "https".
// The host parameter specifies the hostname or IP address of the server to send the request to.
// The baseUri parameter specifies the base URI for the API endpoint, such as "/v1" or "/api".
// The endpoint parameter specifies the endpoint for the API request, such as "/users" or "/posts/1".
//
// Note: Support to use "-" or "" to ignore this parameter.
//
// Example usage:
//
//	client.SetURL("https", "example.com", "/api/v1", "/users")
func (c *Client[T]) SetURL(scheme, host, baseUri, endpoint string) *Client[T] {
	// Set each url component
	c.SetSchema(scheme)
	c.SetHost(host)
	c.SetBaseURI(baseUri)
	c.SetEndpoint(endpoint)

	return c
}

// SetSchema sets the protocol scheme for the client instance.
//
// This method is called by the SetURL method to set the complete URL for the client instance.
//
// See SetURL.
func (c *Client[T]) SetSchema(scheme string) *Client[T] {
	if isEmptyString(scheme) {
		scheme = ProtocolHttp
	}

	if scheme != ProtocolHttp && scheme != ProtocolHttps {
		c.ChalkStr(LogLevelPanic, "scheme parameter: only support http/https protocol header.")
		panic("scheme parameter: only support http/https protocol header")
	}

	c.urls.scheme = scheme

	return c
}

// SetHost sets the host URL for the client instance.
//
// This method is called by the SetURL method to set the complete URL for the client instance.
//
// See SetURL.
func (c *Client[T]) SetHost(host string) *Client[T] {
	if isEmptyString(host) {
		host = fmt.Sprintf("%s:%d", localHost, localPort)
	}

	if !isValidHost(host) && !isValidIPAddrPort(host) {
		c.ChalkStr(LogLevelPanic, "Invalid host or IP address with port.")
		panic("Invalid host or IP address with port.")
	}

	c.urls.host = strings.TrimRight(host, signSlash)

	return c
}

// SetBaseURI sets the base URI for the client instance.
//
// This method is called by the SetURL method to set the complete URL for the client instance.
//
// See SetURL.
func (c *Client[T]) SetBaseURI(baseUri string) *Client[T] {
	if isEmptyString(baseUri) && c.Config.IsDebug {
		c.ChalkStr(LogLevelDebug, "BaseUri parameter is not set, it is recommended to follow the principle of minimal URLs when setting it.")
	}

	signCollect := [2]string{signSlash, signHorizontal}

	for _, v := range signCollect {
		c.urls.baseURI = strings.TrimRight(baseUri, v)
	}

	return c
}

// SetEndpoint sets the API endpoint for the client instance to the specified endpoint.
//
// This method is called by the SetURL method to set the complete URL for the client instance.
//
// See SetURL.
func (c *Client[T]) SetEndpoint(endpoint string) *Client[T] {
	if isEmptyString(endpoint) {
		if c.Config.IsDebug {
			c.ChalkStr(LogLevelDebug, "No access URL endpoint is set, the root path / will be accessed by default.")
		}
		c.urls.endpoint = RootURL
	} else {
		if c.Config.FilterSlash {
			// Turn on the option to automatically mask the trailing /
			c.urls.endpoint = strings.TrimRight(endpoint, signSlash)
		} else {
			// keep original url
			c.urls.endpoint = endpoint
		}
	}

	return c
}

// SetRequest sets the request method and path for the client.
// It also supports dynamic routing by replacing path parameters with actual values.
// The method parameter specifies the HTTP method (e.g., "GET", "POST", "PUT", etc.).
// The path parameter specifies the request path with optional path parameters.
// The pathParams parameter is an optional variadic parameter that contains the values
// to replace the path parameters in the request path.
// It returns the modified client instance.
//
// Dynamic Routing:
// The SetRequest function supports dynamic routing by allowing you to replace path
// parameters in the request path with actual values. If there is only one path
// parameter, you can use the placeholder ":id" in the path. If there are two path
// parameters, the first one should be replaced with ":id" and the second one with ":sid".
//
// Example:
//
//	client := NewClient()
//	client.SetRequest("GET", "/users/:id", "123")
//	client.SetRequest("POST", "/users/:id/:sid", "123", "456")
//
// In the above example, the first SetRequest call sets the request method to "GET" and
// the request path to "/users/123". The second SetRequest call sets the request method
// to "POST" and the request path to "/users/123/456".
//
// Note:
// The SetRequest function currently supports a maximum of two dynamic routing parameters.
// If more than two path parameters are provided, a panic will occur.
func (c *Client[T]) SetRequest(method, path string, pathParams ...string) *Client[T] {
	// Parse Dynamic Routing
	var tempPath string
	switch len(pathParams) {
	case 0:
		tempPath = path
	case 1:
		tempPath = strings.Replace(path, ":id", pathParams[0], 1)
	case 2:
		tempPath = strings.Replace(path, ":id", pathParams[0], 1)
		tempPath = strings.Replace(tempPath, ":sid", pathParams[1], 1)
	case 3:
		panic("There are too many dynamic routing parameters, which are not supported for now!")
	}

	// Parse the URL
	parseUrl := urlSegments(tempPath)

	// Set the request method
	c.SetMethod(method)

	// Set the URL
	c.SetURL(parseUrl.scheme, parseUrl.host, parseUrl.baseURI, parseUrl.endpoint)

	// Set the query parameters
	if method != MethodOptions && !isEmpty(parseUrl.params) {
		c.SetQueryParams(parseUrl.params)
	}

	return c
}

// SetQueryParam sets a query parameter for the request.
// It takes a key and value as parameters and adds them to the params map of the Client instance.
// It returns a pointer to the Client instance, allowing for method chaining.
//
// Example usage:
//
//	client.SetQueryParam("key", "value")
func (c *Client[T]) SetQueryParam(key, value string) *Client[T] {
	c.params[key] = value

	return c
}

// SetQueryParams sets multiple query parameters for the request.
// It takes a `params` map as a parameter and converts it to the `SMap` type, which is used to
// store query parameters in the `Client` instance.
// The `params` map contains key-value pairs where the keys represent the query parameter names
// and the values represent the query parameter values.
// This method replaces any existing query parameters in the `params` map of the `Client`
// instance with the new ones.
// It returns a pointer to the `Client` instance to allow for method chaining.
//
// Example usage:
//
//	params := H{"key1": "value1", "key2": "value2"}
//	client.SetQueryParams(params)
func (c *Client[T]) SetQueryParams(params H) *Client[T] {
	tempParams := convertToSMap(params)

	if isEmpty(c.params) {
		c.params = tempParams
		return c
	}

	for key, value := range tempParams {
		c.params[key] = value
	}
	return c
}

// SetHeader sets a custom header for the request.
// It takes a `key` and `value` as parameters and adds the header to the `Client` instance.
// The `key` parameter represents the header key, and the `value` parameter represents the header value.
// This method allows adding custom headers to the request.
// It returns a pointer to the `Client` instance to allow for method chaining.
//
// Example usage:
//
//	client.SetHeader("Content-Type", "application/json")
func (c *Client[T]) SetHeader(key, value string) *Client[T] {
	c.headers.extra[key] = value

	return c
}

// SetHeaders sets multiple custom headers for the request.
// It takes a `headers` parameter, which is a map[string]string representing the headers to be set.
// Each key-value pair in the `headers` map corresponds to a header key and value, respectively.
// This method allows setting multiple custom headers at once.
// It returns a pointer to the `Client` instance to allow for method chaining.
//
// Example usage:
//
//	headers := map[string]any{
//		"Content-Type":  "application/json",
//		"Authorization": "Bearer <token>",
//	}
//	client.SetHeaders(headers)
func (c *Client[T]) SetHeaders(headers H) *Client[T] {
	parseHeaders := convertToSMap(headers)

	if isEmpty(c.headers.extra) {
		c.headers.extra = parseHeaders
		return c
	}

	for key, value := range parseHeaders {
		c.headers.extra[key] = value
	}
	return c
}

// SetCookie adds a cookie to the request headers.
// It takes a `cookie` parameter, which is a pointer to an `http.Cookie` representing the cookie to be added.
// This method allows adding a cookie to the request headers.
// It returns a pointer to the `Client` instance to allow for method chaining.
//
// Example usage:
//
//	cookie := &http.Cookie{
//		Name:  "session",
//		Value: "1234567890",
//	}
//	client.SetCookie(cookie)
func (c *Client[T]) SetCookie(cookie *http.Cookie) *Client[T] {
	c.headers.cookies = append(c.headers.cookies, cookie)

	return c
}

// SetCookies sets the cookies for the request headers.
// It takes a `cookies` parameter, which is a slice of pointers to `http.Cookie` representing the cookies to be set.
// This method allows setting multiple cookies for the request headers.
// It returns a pointer to the `Client` instance to allow for method chaining.
//
// Example usage:
//
//	cookie1 := &http.Cookie{
//		Name:  "session",
//		Value: "1234567890",
//	}
//	cookie2 := &http.Cookie{
//		Name:  "user",
//		Value: "john.doe",
//	}
//	cookies := []*http.Cookie{cookie1, cookie2}
//	client.SetCookies(cookies)
func (c *Client[T]) SetCookies(cookies []*http.Cookie) *Client[T] {
	c.headers.cookies = cookies

	return c
}

// SetBasicAuth sets the Basic Authentication credentials for the request.
// It takes a `username` and `password` as parameters and sets them as the Basic Authentication
// credentials in the `Client` instance.
// The `username` parameter represents the username for the Basic Authentication.
// The `password` parameter represents the password for the Basic Authentication.
// This method replaces any existing authorization credentials in the `Client` instance with the
// new Basic Authentication credentials.
// It returns a pointer to the `Client` instance to allow for method chaining.
//
// Example usage:
//
//	client.SetBasicAuth("username", "password")
func (c *Client[T]) SetBasicAuth(username, password string) *Client[T] {
	c.authorization = &authorization{
		authType:      AuthTypeBasic,
		basicUsername: username,
		basicPassword: password,
		bearerToken:   "",
	}

	return c
}

// SetBearerAuth sets the Bearer Token for the request.
// It takes a `token` as a parameter and sets it as the Bearer Token in the `Client` instance.
// The `token` parameter represents the Bearer Token for the authentication.
// This method replaces any existing authorization credentials in the `Client` instance with
// the new Bearer Token.
// It returns a pointer to the `Client` instance to allow for method chaining.
//
// Example usage:
//
//	client.SetBearerAuth("your-token")
func (c *Client[T]) SetBearerAuth(token string) *Client[T] {
	c.authorization = &authorization{
		authType:      AuthTypeBearer,
		bearerToken:   token,
		basicUsername: "",
		basicPassword: "",
	}

	return c
}

// SetAccept sets the value of the "Accept" header for the request.
// It takes an `accept` parameter, which is a string representing the value of the "Accept" header.
// This method allows specifying the desired media type for the response.
// It returns a pointer to the `Client` instance to allow for method chaining.
//
// Example usage:
//
//	client.SetAccept("application/json")
func (c *Client[T]) SetAccept(accept string) *Client[T] {
	c.headers.accept = accept

	return c
}

// SetContentType sets the value of the "Content-Type" header for the request.
// It takes a `ct` parameter, which is a string representing the value of the "Content-Type" header.
// This method allows specifying the media type of the request payload.
// It returns a pointer to the `Client` instance to allow for method chaining.
//
// Example usage:
//
//	client.SetContentType("application/json")
func (c *Client[T]) SetContentType(ct string) *Client[T] {
	c.headers.contentType = ct

	return c
}

// SetLanguage sets the value of the "Accept-Language" header for the request.
// It takes a `lang` parameter, which is a string representing the value of the "Accept-Language" header.
// This method allows specifying the preferred language for the response.
// It returns a pointer to the `Client` instance to allow for method chaining.
//
// Example usage:
//
//	client.SetLanguage("en-US")
func (c *Client[T]) SetLanguage(lang string) *Client[T] {
	c.headers.language = lang

	return c
}

// SetUserAgent sets the value of the "User-Agent" header for the request.
// It takes a `ua` parameter, which is a string representing the value of the "User-Agent" header.
// This method allows specifying the user agent for the request.
// It returns a pointer to the `Client` instance to allow for method chaining.
//
// Example usage:
//
//	client.SetUserAgent("MyApp/1.0")
func (c *Client[T]) SetUserAgent(ua string) *Client[T] {
	c.headers.userAgent = ua

	return c
}

// SetJsonPayload sets the JSON payload for the request.
// It takes a `data` parameter, which is a map[string]any representing the JSON data to be sent in the request body.
// This method is used for making JSON-encoded POST or PUT requests.
// It returns a pointer to the `Client` instance to allow for method chaining.
//
// Example usage:
//
//	payload := map[string]interface{}{
//		"name":  "John Doe",
//		"email": "john.doe@example.com",
//	}
//	client.SetJsonPayload(payload)
func (c *Client[T]) SetJsonPayload(data H) *Client[T] {
	c.payload = data

	return c
}

// SetPayload sets the payload for the request.
// It takes a `data` parameter of any type representing the data to be sent in the request body.
// This method is used for making generic POST or PUT requests.
// It returns a pointer to the `Client` instance to allow for method chaining.
//
// Example usage:
//
//	payload := "Hello, World!"
//	client.SetPayload(payload)
func (c *Client[T]) SetPayload(data any) *Client[T] {
	c.payload = data

	return c
}

/*
	Internal chain methods with Setter attribute for the Client struct
*/

// Deprecated: The sonic method in the Client struct is designed to improve performance by optimizing the
// handling of header information. By calling this method before setting additional headers
// using the SetHeader method, significant performance improvements can be achieved.
//
// The sonic method internally utilizes the UsePreHooks method to perform a pre-hook operation.
// In this case, it sets the client.headers to a new header object, which effectively clears
// any existing header information.
//
// By doing so, the method avoids potential performance issues caused by memory allocations and
// extra operations when setting additional headers. This optimization ensures that unnecessary
// header modifications are avoided and reduces potential performance overhead.
//
// By utilizing the sonic method and then setting headers using the SetHeader method, the
// performance of the code is significantly improved. In fact, this approach has demonstrated a
// remarkable 10-fold increase in performance compared to setting headers directly without
// utilizing the sonic method.
//
// Overall, this approach helps enhance performance by optimizing the handling of header
// information in the Client struct, resulting in a remarkable 10-fold performance improvement.
func (c *Client[T]) sonic() *Client[T] {
	c.UsePreHooks(func(client *Client[T]) error {
		client.headers = &header{}
		return nil
	})
	return c
}

/*
	Internal Methods and Functions for Construct Client Struct
*/

// createRequest creates and prepares an HTTP request based on the client instance's configuration.
// It sets the request method, URL, body, headers, authentication, cookies, and other request configurations.
// The created request is stored in the client's context.
//
// This internal function is called by the Send method to set the complete request for the client instance.
//
// See Send.
func (c *Client[T]) createRequest() *Client[T] {
	// Check if request method is set
	if isEmptyString(c.Meta.Method) {
		panic("An incomplete request, must set the request method.")
	}

	// Parsing the full url path and query params
	c.parseFullURLPath()

	// Parsing the request body
	var req *http.Request
	var err error

	// Set request body
	if isEmpty(c.payload) {
		// such as GET
		req, err = http.NewRequest(c.Meta.Method, c.Meta.Url, nil)
	} else {
		// such as POST/PUT...
		var byteData []byte
		byteData, err = c.Config.JSONLoader.Marshal(c.payload)
		if err != nil {
			c.Exception = &Exception{
				CodeLocation:   fileLocation(1),
				PanicError:     err,
				OccurrenceTime: time.Now().Unix(),
			}
			return c
		}
		bodyReader := bytes.NewReader(byteData)
		req, err = http.NewRequest(c.Meta.Method, c.Meta.Url, bodyReader)
	}

	// Store the request object to the context
	c.Context.Request = req
	if err != nil {
		c.Exception = &Exception{
			CodeLocation:   fileLocation(1),
			PanicError:     err,
			OccurrenceTime: time.Now().Unix(),
		}
		return c
	}

	// Set custom request headers
	if len(c.headers.extra) > 0 {
		extraHeaders := make(http.Header, len(c.headers.extra))
		for k, v := range c.headers.extra {
			extraHeaders[k] = []string{v}
		}
		req.Header = extraHeaders
	}

	// Set User-Agent request headers
	if !isEmpty(c.headers.userAgent) {
		req.Header.Set(HeaderUserAgentKey, c.headers.userAgent)
	}

	// Set Accept request headers
	if !isEmpty(c.headers.accept) {
		req.Header.Set(HeaderAcceptKey, c.headers.accept)
	}

	// Set Content-Type request headers
	if !isEmpty(c.headers.contentType) {
		req.Header.Set(HeaderContentTypeKey, c.headers.contentType)
	}

	// Set Content-Language request headers
	if !isEmpty(c.headers.language) {
		req.Header.Set(HeaderContentLanguageKey, c.headers.language)
	}

	// Set Authorization request headers
	switch c.authorization.authType {
	case AuthTypeBasic:
		req.Header.Set(HeaderAuthorizationKey, getBasicAuth(c.authorization.basicUsername, c.authorization.basicPassword))
	case AuthTypeBearer:
		req.Header.Set(HeaderAuthorizationKey, getBearerAuth(c.authorization.bearerToken))
	default:
		// pass
	}

	// Set Cookies request headers
	if !isEmpty(c.headers.cookies) {
		for _, v := range c.headers.cookies {
			req.AddCookie(v)
		}
	}

	// Set client request configs
	client := httpClientDefaultConf(c.Config.Timeout, c.Config.SkipTLS, c.Config.Logger)

	// Store the client object to the context
	c.Context.HttpClient = client

	return c
}

// parseFullURLPath generates the complete URL path for the client instance by concatenating the individual URL components,
// including scheme, host, base URI, endpoint and query parameters.
//
// This internal function is called by the createRequest method to set the complete URL and query param section for the client instance.
//
// See createRequest.
func (c *Client[T]) parseFullURLPath() {
	var urlPath string

	// Set the url path part
	if isEmptyString(c.Meta.Url) {
		u := c.urls
		if u.baseURI == RootURL {
			urlPath = fmt.Sprintf("%s://%s%s%s", u.scheme, u.host, "", u.endpoint)
		} else {
			urlPath = fmt.Sprintf("%s://%s%s%s", u.scheme, u.host, u.baseURI, u.endpoint)
		}
	}

	// Set request parameters section
	switch len(c.params) {
	case 0:
		c.Meta.Url = urlPath
	default:
		// Use url.Values to store query parameters
		queryParams := url.Values{}
		for k, v := range c.params {
			queryParams.Add(k, v)
		}

		// Encode query parameters as URL strings
		encodedQueryParams := queryParams.Encode()

		// Generate the full request path
		fullURL := fmt.Sprintf("%s?%s", urlPath, encodedQueryParams)

		c.Meta.Url = fullURL
	}
}

// httpClientDefaultConf creates and returns a default HTTP client with the specified configurations.
// The timeout parameter specifies the maximum amount of time to wait for a response.
// The skipTLS parameter indicates whether to skip TLS certificate verification.
// The logFmt parameter is an optional logger to log HTTP requests and responses.
func httpClientDefaultConf(timeout time.Duration, skipTLS bool, logFmt *log.Logger) *http.Client {
	// Create a new transport object with the following configurations:
	tr := &http.Transport{
		// TLSClientConfig is set to skip certificate verification.
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: skipTLS,
		},
		// MaxIdleConns specifies the maximum number of idle (keep-alive) connections across all hosts.
		MaxIdleConns: 10,
		// MaxIdleConnsPerHost specifies the maximum number of idle (keep-alive) connections per host.
		MaxIdleConnsPerHost: 10,
		// IdleConnTimeout specifies the maximum amount of time a connection may remain idle (keep-alive)
		// before it is closed and removed from the pool.
		IdleConnTimeout: 60 * time.Second,
	}

	// Create an HTTP client with a timeout for receiving a response.
	client := &http.Client{
		// The maximum amount of time to wait for a response is specified by the Timeout field.
		Timeout: timeout,
		// Use the origin default transport object.
		Transport: http.DefaultTransport,
	}

	if isEmpty(logFmt) {
		// Set the transport object to be used for the HTTP client.
		client.Transport = tr
	} else {
		// Create a custom Logger transport object.
		client.Transport = &loggedTransport{
			transport: tr,
			logger:    logFmt,
		}
	}

	return client
}
