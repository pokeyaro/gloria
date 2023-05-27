// Copyright (c) 2023 Pokeya Boa <pokeya.mystic@gmail.com>, All rights reserved.
// resty source code and usage is governed by a MIT style
// license that can be found in the LICENSE file.

package gloria

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strings"
	"time"
)

const (
	// MethodGet HTTP method
	MethodGet = "GET"

	// MethodPost HTTP method
	MethodPost = "POST"

	// MethodPut HTTP method
	MethodPut = "PUT"

	// MethodDelete HTTP method
	MethodDelete = "DELETE"

	// MethodPatch HTTP method
	MethodPatch = "PATCH"

	// MethodHead HTTP method
	MethodHead = "HEAD"

	// MethodOptions HTTP method
	MethodOptions = "OPTIONS"
)

const (
	// ProtocolHttp HTTP protocol
	ProtocolHttp = "http"

	// ProtocolHttps HTTPS protocol
	ProtocolHttps = "https"
)

const (
	// Localhost IP address
	localHost = "127.0.0.1"

	// Localhost port number
	localPort = 8080
)

const (
	// Slash character "/"
	signSlash = "/"

	// Horizontal line character "-"
	signHorizontal = "-"
)

const (
	// RootURL an alias for signSlash
	RootURL = signSlash

	// Placeholder an alias for signHorizontal
	Placeholder = signHorizontal
)

const (
	// OkCode Predefined default success response codes
	OkCode = 0

	// FailCode Predefine default failure response codes
	FailCode = 50000
)

const (
	// TimeoutShort Short timeout duration (10 seconds)
	TimeoutShort = 10 * time.Second

	// TimeoutMedium Medium timeout duration (30 seconds)
	TimeoutMedium = 30 * time.Second

	// TimeoutLong Long timeout duration (60 seconds)
	TimeoutLong = 60 * time.Second
)

const (
	// AuthTypeBasic Basic authentication type
	AuthTypeBasic = "Basic"

	// AuthTypeBearer Bearer authentication type
	AuthTypeBearer = "Bearer"
)

const (
	// LocaleEn English locale
	LocaleEn = "en-US,en;q=0.9"

	// LocaleZh Chinese locale
	LocaleZh = "zh-CN,zh;q=0.9"
)

const (
	// PlainTextType Plain text content type
	PlainTextType = "text/plain; charset=utf-8"

	// JsonContentType JSON content type
	JsonContentType = "application/json"

	// FormContentType Form data content type
	FormContentType = "application/x-www-form-urlencoded"
)

var (
	// QueryMethods Method list in a String Slice
	QueryMethods = []string{
		MethodGet,
		MethodPost,
		MethodPut,
		MethodDelete,
		MethodPatch,
		MethodHead,
		MethodOptions,
	}
)

var (
	HeaderAcceptKey          = http.CanonicalHeaderKey("Accept")
	HeaderLocationKey        = http.CanonicalHeaderKey("Location")
	HeaderUserAgentKey       = http.CanonicalHeaderKey("User-Agent")
	HeaderContentTypeKey     = http.CanonicalHeaderKey("Content-Type")
	HeaderContentLengthKey   = http.CanonicalHeaderKey("Content-Length")
	HeaderContentLanguageKey = http.CanonicalHeaderKey("Content-Language")
	HeaderContentEncodingKey = http.CanonicalHeaderKey("Content-Encoding")
	HeaderAuthorizationKey   = http.CanonicalHeaderKey("Authorization")
)

type Client[T any] struct {
	// context
	Context *Context

	// request metadata
	Meta *Meta

	// request settings
	Config *Config

	// intercepted error
	Exception *Exception

	// expandable body
	Result *RESTFulResp[T]

	// middlewares
	beforeRequest []func(*Client[T]) error
	afterResponse []func(*Client[T]) error

	// request content
	urls          *urls
	params        SMap
	authorization *authorization
	headers       *header
	payload       any
}

// H is a type alias for an exported map[string]interface{}
type H = map[string]any

// SMap is a type alias for an exported map[string]string
type SMap = map[string]string

// timestamp is a type alias for int64
type timestamp = int64

// Request urls related
type urls struct {
	scheme   string
	host     string
	baseURI  string
	endpoint string
}

// rawUrl represents a URL with additional fields to handle query parameters
type rawUrl struct {
	urls
	params H
}

// Request authorization related
type authorization struct {
	authType      string
	bearerToken   string
	basicUsername string
	basicPassword string
}

// Request header related
type header struct {
	accept      string
	contentType string
	language    string
	cookies     []*http.Cookie
	userAgent   string
	extra       SMap
}

type Context struct {
	// client object
	HttpClient *http.Client

	// raw request
	Request *http.Request

	// raw response
	Response *Response
}

type Meta struct {
	Method     string        // store the request method
	Url        string        // store the full url path
	Duration   time.Duration // time-consuming current request
	ReceivedAt time.Time     // store the timestamp indicating when the response was received
}

type Config struct {
	Timeout       time.Duration
	SkipTLS       bool
	FilterSlash   bool
	IsDebug       bool
	Logger        *log.Logger
	IsRestMode    bool
	DefaultOkCode int
	JSONLoader    JSONLibrary
}

type Exception struct {
	CodeLocation   string
	PanicError     error
	FailureReason  string
	OccurrenceTime timestamp
}

type RESTFulResp[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data,omitempty"`
}

type Response struct {
	R      *http.Response // all native objects
	Status int            // http response status code

	bs     []byte
	text   string
	length int64
}

func (c *Client[T]) Send() *Client[T] {
	// request middleware
	for _, md := range c.beforeRequest {
		if err := md(c); err != nil {
			c.Exception = &Exception{
				CodeLocation:   fileLocation(1),
				PanicError:     err,
				OccurrenceTime: time.Now().Unix(),
			}
			return c
		}
	}

	// create
	c.createRequest()

	// record start time
	startTime := time.Now()

	// execute
	resp, err := c.Context.HttpClient.Do(c.Context.Request)

	if err != nil {
		c.Exception = &Exception{
			CodeLocation:   fileLocation(1),
			PanicError:     err,
			OccurrenceTime: time.Now().Unix(),
		}
		return c
	}

	defer func() {
		if err = resp.Body.Close(); err != nil {
			// Handle Close() errors, such as logging or returning an error message
			c.Exception = &Exception{
				CodeLocation:   fileLocation(1),
				PanicError:     err,
				OccurrenceTime: time.Now().Unix(),
			}
			c.ChalkPrintf(LogLevelPanic, "Failed to close response body:", err)
		}
	}()

	// record end time
	duration := time.Since(startTime)
	c.Meta.Duration = duration

	// record received At
	c.Meta.ReceivedAt = time.Now()

	// response middleware
	for _, md := range c.afterResponse {
		if err = md(c); err != nil {
			c.Exception = &Exception{
				CodeLocation:   fileLocation(1),
				PanicError:     err,
				OccurrenceTime: time.Now().Unix(),
			}
			return c
		}
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.Exception = &Exception{
			CodeLocation:   fileLocation(1),
			PanicError:     err,
			OccurrenceTime: time.Now().Unix(),
		}
		return c
	}

	c.Context.Response = &Response{
		R:      resp,
		Status: resp.StatusCode,
		bs:     body,
		text:   string(body),
		length: resp.ContentLength,
	}

	if c.Context.Response.length == 0 {
		c.Exception = &Exception{
			CodeLocation:   fileLocation(1),
			PanicError:     errors.New("response body length is 0"),
			OccurrenceTime: time.Now().Unix(),
		}
		return c
	}

	var errJson error
	if c.Config.IsRestMode {
		errJson = c.Config.JSONLoader.Unmarshal(c.Context.Response.bs, &c.Result)
	} else {
		errJson = c.Config.JSONLoader.Unmarshal(c.Context.Response.bs, &c.Result.Data)
	}
	if errJson != nil {
		c.Exception = &Exception{
			CodeLocation:   fileLocation(1),
			PanicError:     errJson,
			OccurrenceTime: time.Now().Unix(),
		}
		return c
	}

	if c.Config.IsDebug {
		c.ChalkStr(LogLevelDebug, c.Context.Response.text)
	}

	if c.Context.Response.Status != http.StatusOK {
		c.Exception = &Exception{
			CodeLocation:   fileLocation(1),
			FailureReason:  c.Result.Msg,
			OccurrenceTime: time.Now().Unix(),
		}
	}

	return c
}

func (c *Client[T]) Unwrap() (*Client[T], string) {
	if c.Exception.PanicError != nil {
		panic(c.Exception.PanicError.Error())
	}
	if c.Exception.FailureReason != "" {
		return c, fmt.Sprintf(
			`HTTP request method: [%s], HTTP request url path: "%s", HTTP response status code and description: "%s", business error code: %d, business error reason: "%s", occurrence time: %v\n`,
			c.Meta.Method,
			c.Meta.Url,
			c.Context.Response.R.Status,
			c.Result.Code,
			c.Exception.FailureReason,
			c.Exception.OccurrenceTime,
		)
	}
	return c, ""
}

func (c *Client[T]) Data() T {
	return c.Result.Data
}

func (c *Client[T]) EchoQPS() float64 {
	seconds := c.Meta.Duration.Seconds()
	qps := float64(1) / seconds
	if c.Config.IsDebug {
		c.ChalkPrintf(LogLevelDebug, "An approximate calculation of Queries Per Second (QPS) yields a result of: %.6f. Please note that this calculation may not be entirely accurate.", qps)
	}
	return qps
}

func (c *Client[T]) EchoTime() (time.Duration, time.Time) {
	return c.Meta.Duration, c.Meta.ReceivedAt
}

func (c *Client[T]) EchoBenchmark() (int, int64) {
	count := int(math.Round(c.EchoQPS()))
	nanoseconds := c.Meta.Duration.Nanoseconds()
	return count, nanoseconds
}

func (c *Client[T]) EchoProto() string {
	return c.Context.Response.R.Proto
}

func (c *Client[T]) EchoCode() (int, int) {
	httpStatusCode := c.Context.Response.R.StatusCode
	restReturnCode := c.Result.Code
	return httpStatusCode, restReturnCode
}

func (c *Client[T]) EchoMessage() (string, string) {
	httpStatusMsg := strings.Join(strings.Split(c.Context.Response.R.Status, " ")[1:], "")
	restReturnMsg := c.Result.Msg
	return httpStatusMsg, restReturnMsg
}

func (c *Client[T]) EchoMode() string {
	if c.Config.IsRestMode {
		return "RESTFul Response"
	}
	return "HTTP Response"
}

func (c *Client[T]) EchoURL() (string, string) {
	return c.Meta.Method, c.Meta.Url
}

func (c *Client[T]) Echo() {
	var output strings.Builder
	qps := c.EchoQPS()
	executions, efficiency := c.EchoBenchmark()
	proto := c.EchoProto()
	method, fullpath := c.EchoURL()
	statusCode, errCode := c.EchoCode()
	statusMsg, errMsg := c.EchoMessage()
	durationTime, receivedAt := c.EchoTime()
	mode := c.EchoMode()
	output.WriteString("[API Call Insights]\n")
	if !isEmpty(c.Exception.FailureReason) && !isEmpty(c.Exception.PanicError) {
		output.WriteString(fmt.Sprintf("  Panic      : %v\n", c.Exception.PanicError))
		output.WriteString(fmt.Sprintf("  Reason     : %s\n", c.Exception.FailureReason))
		output.WriteString(fmt.Sprintf("  Location   : %s\n", c.Exception.CodeLocation))
		output.WriteString(fmt.Sprintf("  Occurrence : %v\n", c.Exception.OccurrenceTime))
	} else {
		output.WriteString(fmt.Sprintf("  Mode       : %s\n", mode))
		output.WriteString(fmt.Sprintf("  Error      : %v\n", (*any)(nil)))
		output.WriteString(fmt.Sprintf("  Method     : %s\n", method))
		output.WriteString(fmt.Sprintf("  URL        : %s\n", fullpath))
		if !c.Config.IsRestMode {
			output.WriteString(fmt.Sprintf("  Status     : %s\n", c.Context.Response.R.Status))
		} else {
			output.WriteString(fmt.Sprintf("  Status Code: %d\n", statusCode))
			output.WriteString(fmt.Sprintf("  Status Desc: %s\n", statusMsg))
			output.WriteString(fmt.Sprintf("  Return Code: %d\n", errCode))
			output.WriteString(fmt.Sprintf("  Return Msg : %s\n", errMsg))
		}
		output.WriteString(fmt.Sprintf("  Benchmark  : %d\t%d ns/op\n", executions, efficiency))
		output.WriteString(fmt.Sprintf("  Proto      : %s\n", proto))
		output.WriteString(fmt.Sprintf("  QPS        : %.6f\n", qps))
		output.WriteString(fmt.Sprintf("  Duration   : %v\n", durationTime))
		output.WriteString(fmt.Sprintf("  Received At: %s\n", receivedAt.Format(time.RFC850)))
		output.WriteString(fmt.Sprintf("  Body       : %v\n", "-"))
	}
	fmt.Println(output.String())
}

func (c *Client[T]) ToJson(v any) error {
	if c.Context.Response.length == 0 {
		return errors.New("pesponse body length is 0")
	}
	if err := c.Config.JSONLoader.Unmarshal(c.Context.Response.bs, v); err != nil {
		return err
	}
	return nil
}

type beforeRequest[T any] func(*Client[T]) error
type afterResponse[T any] func(*Client[T]) error

// UsePreHooks request interceptor middleware
func (c *Client[T]) UsePreHooks(funcs ...beforeRequest[T]) {
	if c.Config.IsDebug {
		c.ChalkStr(LogLevelDebug, "inject pre hooks")
	}
	for _, fn := range funcs {
		c.beforeRequest = append(c.beforeRequest, fn)
	}
}

// UsePostHooks response Interceptor Middleware
func (c *Client[T]) UsePostHooks(funcs ...afterResponse[T]) {
	if c.Config.IsDebug {
		c.ChalkStr(LogLevelDebug, "inject post hooks")
	}
	for _, fn := range funcs {
		c.afterResponse = append(c.afterResponse, fn)
	}
}
