<div align=center>
  <img src="logo.png" width="450" height="225" alt="gloria" />

  <br/>

  ![Go version](https://img.shields.io/badge/go-%3E%3Dv1.18-9cf)
  ![Release](https://img.shields.io/badge/release-1.0.0-green.svg)
  [![GoDoc](https://godoc.org/github.com/pokeyaro/gloria?status.svg)](https://godoc.org/github.com/pokeyaro/gloria)
  [![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
</div>

<p style="font-size: 15px">
  Gloria is a minimalist and elegant HTTP client tool designed to embrace the beauty of RESTful style. 
  Gloria draws inspiration from 
  <a href="https://github.com/go-resty">Go resty</a>, 
  <a href="https://github.com/requests">Python requests</a>, 
  and <a href="https://github.com/axios">JavaScript axios</a>.
</p>

English | [简体中文](./README_zh.md)


## Origin

The term `"Gloria"` is derived from the [hit song](https://www.youtube.com/watch?v=stGUpMav1sc) in G.E.M.'s album "Revelation". We hope that this library can provide a wonderful experience just like G.E.M.'s music.

Additionally, `Gloria` carries the meaning of "glory" in Latin, symbolizing the library's design to offer a glorious and elegant way of using the `Go` `HTTP` client. Similar to the `RESTful` specification, it aims to make network requests more in line with human natural habits.

We pay tribute to all those who strive to build internet applications and hope that `Gloria` can bring them a more comfortable development and user experience.


## Features

- 🪶 Simple, user-friendly, and versatile `API` design
- 👏 Unified `RESTful-style` response implementation using generics
- 🚀 Support for injecting multiple request and response interceptors
- 📝 Detailed colored logging and error tracing capabilities
- 💃🏼 Elegant error handling with `Unwrap` (`Rust-style`)
- 🧭 Request invocation time and `QPS` estimation functionality
- 🌈 Support for `GET`, `POST`, `PUT`, `DELETE`, and more
- 🎋 Thorough `Benchmark` testing for optimal performance


## Usage

### Version

Supported versions: Requires the `Go` `T` generics feature and a minimum interpreter version of `go1.18+`.

### Installation

```bash
go get -u github.com/pokeyaro/gloria
```

### Import

```go
import "github.com/pokeyaro/gloria"
```


## Index

- [Simple GET Request](#SimpleGet)
- [More Advanced Request Configuration](#ExtendedConf)
- [Common CRUD Request Operations](#CrudMethod)
- [Two Response Modes](#TwoRespMode)
- [REST Syntactic Sugar](#RestSyntacticSugar)
- [More API Function Signatures](#MoreAPI)
- [Injecting Interceptors for HTTP](#HttpInterceptor)
- [Some Coding Suggestions](#CodeSuggestions)


## Documentation

Through the following examples, you can quickly get started and learn about the `gloria` library.

### <span id="SimpleGet">Simple GET Request</span>

> For a real API, you can use: `curl -X GET --location 'http://httpbin.org/get'`

#### Example 1

To demonstrate how to use `gloria` in a clear way, we will use the `NewHTTP()` method and explain the steps involved:

Example code: [api-httpbin.go](./example_test.go) &nbsp; | &nbsp; Recommended 🌟🌟🌟

```go
// Preparing the response body struct.
type HttpBin struct {
    Args struct {
    } `json:"args"`
    Headers struct {
        Accept                  string `json:"Accept"`
        AcceptEncoding          string `json:"Accept-Encoding"`
		AcceptLanguage          string `json:"Accept-Language"`
        Host                    string `json:"Host"`
        UpgradeInsecureRequests string `json:"Upgrade-Insecure-Requests"`
        UserAgent               string `json:"User-Agent"`
        XAmznTraceId            string `json:"X-Amzn-Trace-Id"`
    } `json:"headers"`
    Origin string `json:"origin"`
    Url    string `json:"url"`
}

// Building a request using the NewHTTP method.
client := gloria.NewHTTP[HttpBin]()

// Setting the request type.
client.SetMethod(gloria.MethodGet)

// Setting the request's request route resource (URL segments: proto, host, baseURI, endpoint).
client.SetURL(gloria.ProtocolHttp, "httpbin.org", "", "/get")

// Sending the request.
client.Send()

// Handling errors.
client.Unwrap()

// Printing request metadata.
client.Echo()

/* Output:
[API Call Insights]
  Mode       : HTTP Response
  Error      : <nil>
  Method     : GET
  URL        : http://httpbin.org/get
  Status     : 200 OK
  Benchmark  : 1	1534950042 ns/op
  Proto      : HTTP/1.1
  QPS        : 0.651487
  Duration   : 1.534950042s
  Received At: Saturday, 20-May-23 23:28:11 CST
  Body       : -
*/
```

#### Example 2

Providing users with a more concise `API` method and supporting chain loading, so you can write it like this:

Example code: [api-httpbin.go](./example_test.go) &nbsp; | &nbsp; Recommended 🌟🌟🌟🌟🌟

```go
type HttpBin struct {
    // omitted...
}

// Request
client, _ := gloria.NewHTTP[HttpBin]().SetRequest(gloria.MethodGet, "http://httpbin.org/get").Send().Unwrap()

// Retrieve
fmt.Println(client.Data().Url)

/* Output:
http://httpbin.org/get
*/
```

### <span id="ExtendedConf">More Advanced Request Configuration</span>

> For a real API, you can use: `curl -X GET --location 'https://api.thecatapi.com/v1/images/search?size=med&mime_types=png%2Cgif&format=json&order=RANDOM&limit=20' --header 'Content-Type: application/json'`

#### Example 1

Let's experience the richer configuration options provided by `New()` and `Optional()` methods:

Example code: [api-cat.go](./examples/request_test.go) &nbsp; | &nbsp; Recommended 🌟🌟🌟🌟

```go
// Preparing the response body struct.
type ImageSearch struct {
    Id     string `json:"id"`
    Url    string `json:"url"`
    Width  int    `json:"width"`
    Height int    `json:"height"`
}

// Build a more comprehensive request.
func GetCatAPI() {
	// This API returns a slice of struct.
    r := gloria.New[[]ImageSearch]()

    r.Optional(
        // Whether to enable debug mode.
        gloria.WithIsDebug[[]ImageSearch](false),
        // Whether to enable logging.
        gloria.WithUseLogger[[]ImageSearch](true),
        // Setting additional parameters using the lambda syntax.
        gloria.Lambda[[]ImageSearch](func(c *gloria.Client[[]ImageSearch]) {
            c.Config.SkipTLS = true
            c.Config.Timeout = gloria.TimeoutMedium
            c.Config.IsRestMode = false // If it's a non-standard RESTful response interface, use native HTTP mode by selecting false.
        }),
    ).
        // Set the request method.
        SetMethod(gloria.MethodGet).
        // Set the request path. This method requires specifying the path in segments.
        SetURL(gloria.ProtocolHttps, "api.thecatapi.com", "/v1", "/images/search").
        // Set multiple request parameters.
		SetQueryParams(gloria.H{
            "size":       "med",
            "mime_types": []string{"png", "gif"},
            "format":     "json",
            "order":      "RANDOM",
            "limit":      20,
        }).
        // Set multiple request headers.
        SetHeaders(gloria.H{
            "x-api-key": "live_example-api-key",
            "Content-Type": "application/json",
        }).
        // Sending the request.
        Send().
        Unwrap()

    fmt.Println(r.Data()[0].Url)	
}
```

#### Example 2

Use `Default()` to load default configurations. 
In fact, you don't need to worry about the difference between `Default()` and `New()`, or what default parameters it adds internally. 
If you are curious and want to know more, you can use the `go doc gloria.Default` command to check.

Example code: [api-cat.go](./examples/request_test.go) &nbsp; | &nbsp; Recommended 🌟🌟🌟🌟🌟

```go
type ImageSearch struct {
    // omitted...
}

// Build a more sophisticated request.
func GetCatAPI() {
	// Note that the Default() method operates in REST mode. To switch to another mode, use ToggleMode().
    r := gloria.Default[[]ImageSearch]().ToggleMode()

    r.
        SetRequest(gloria.MethodGet, "https://api.thecatapi.com/v1/images/search").
        SetQueryParams(gloria.H{
            "size":       "med",
            "mime_types": []string{"png", "gif"},
            "format":     "json",
            "order":      "RANDOM",
            "limit":      20,
        }).
        SetHeader("x-api-key", "live_example-api-key").
        SetContentType(gloria.JsonContentType).
        Send().Unwrap()

    fmt.Println(r.Result.Data[0].Url)
}
```

### <span id="CrudMethod">Common CRUD Request Operations</span>

Below demonstrates the `CRUD` operations of the `API`, including `Create` `[POST]`, `Read` `[GET]`, `Update` `[PUT]`, and `Delete` `[DELETE]`.

#### Example of [GET] Request

Example code: [api-cat.go](./examples/request_test.go) &nbsp; | &nbsp; Recommended 🌟🌟🌟🌟🌟

```go
type FavouritesList struct {
    Id        int       `json:"id"`
    UserId    string    `json:"user_id"`
    ImageId   string    `json:"image_id"`
    SubId     string    `json:"sub_id"`
    CreatedAt time.Time `json:"created_at"`
    Image     struct {
        Id  string `json:"id,omitempty"`
        Url string `json:"url,omitempty"`
    } `json:"image"`
}

// GET
func main() {
    r := gloria.NewHTTP[[]FavouritesList]()

    r.SetRequest(gloria.MethodGet, "https://api.thecatapi.com/v1/favourites").SetHeaders(gloria.H{
        "x-api-key":    "your-api-key",
        "Content-Type": "application/json",
    }).Send().Unwrap()

    for _, v := range r.Data() {
        fmt.Println(v)
    }
}
```

#### Example of [POST] Request

Example code: [api-cat.go](./examples/request_test.go) &nbsp; | &nbsp; Recommended 🌟🌟🌟🌟🌟

```go
type FavouriteImgResp struct {
    Message string `json:"message"`
    Id      int    `json:"id"`
}

type FavouriteImgBody struct {
    ImageId string `json:"image_id"`
    SubId   string `json:"sub_id"`
}

// POST
func main() {
    data := FavouriteImgBody{
        ImageId: "12345",
        SubId:   "my-key-12345",
    }

    r := gloria.NewHTTP[FavouriteImgResp]()

    r.SetRequest(gloria.MethodPost, "https://api.thecatapi.com/v1/favourites").SetHeaders(gloria.H{
        "x-api-key":    "your-api-key",
        "Content-Type": "application/json",
    }).SetPayload(&data).Send().Unwrap()

    fmt.Println("post_id:", r.Data().Id)
}
```

#### Example of [PUT] Request

No example available. Please refer to the `POST` method for reference.

#### Example of [DELETE] Request

Example code: [api-cat.go](./examples/request_test.go) &nbsp; | &nbsp; Recommended 🌟🌟🌟🌟🌟

```go
type Result struct {
    Message string `json:"message"`
}

// DELETE
func main() {
    r := gloria.NewHTTP[Result]()

    r.SetRequest(gloria.MethodDelete, "https://api.thecatapi.com/v1/favourites/:id", "232338734").SetHeaders(gloria.H{
        "x-api-key":    "your-api-key",
        "Content-Type": "application/json",
    }).Send().Unwrap()

    fmt.Println("message:", r.Data().Message)
}
```

### <span id="TwoRespMode">Two Response Modes</span>

There are two different response modes: `HTTP` mode and `REST` mode. Let's explain the differences between them!

In the `HTTP` mode, the response data can be of any format. As shown in the previous example, there are no specific conventions or standards to follow. It provides more flexibility.

In the `REST` mode, if your response follows the format shown below:

```json
{
  "code": 0,
  "msg": "success",
  "data": Object,
}
```

Understood! In the `REST` mode, we directly parse the data in the data field, making your `data` request and response more straightforward and direct.

Below is a table summarizing the `API` usage:

<table>
  <tr>
    <th width="15%">Response Format</th>
    <th width="35%">Method Selection</th>
    <th width="auto">Scenario Description</th>
  </tr>
  <tr>
    <td>HTTP mode</td>
    <td>
      <span>Method 1:</span>&nbsp;<code>New().ToggleMode()</code> <br/>
      <span>Method 2:</span>&nbsp;<code>Default().ToggleMode()</code> <br/>
      <span>Method 3:</span>&nbsp;<code>NewHTTP()</code> // Essence: Same as Method 1. <br/>
    </td>
    <td>
      If your response format is non-standard for RESTful response interface, please use the more general HTTP mode.
    </td>
  </tr>
  <tr>
    <td>REST mode</td>
    <td>
      <span>Method 1:</span>&nbsp;<code>New()</code> <br/>
      <span>Method 2:</span>&nbsp;<code>Default()</code> <br/>
      <span>Method 3:</span>&nbsp;<code>NewREST()</code> // Essence: Same as Method 1. <br/>
    </td>
    <td>When your response format follows the standard RESTful response interface, like: 
      <code>{"code": 0, "msg": "success", "data": null}</code>
    </td>
  </tr>
</table>

### <span id="RestSyntacticSugar">REST Syntactic Sugar</span>

> Advantage: The syntax is simpler, and all configurations are parsed and sent within this function. 
> Therefore, it can be seen as an advantage as well as a disadvantage.
>
> Disadvantage: As mentioned above, it does not currently support more advanced syntax, such as chaining 
> other settings directly within these function calls. Hence, it limits the flexibility to inject other 
> configurations in a chain-like manner.

When using the **REST** mode to construct and handle requests, two additional convenient syntax sugars are provided.

#### Python Requests-style

Example code: [requests-style.go](./examples/requests-style/main.go) &nbsp; | &nbsp; Recommended 🌟🌟🌟🌟🌟

```textmate
func Request[T any](path string, params H, data any, headers ...H) ExecMethod[T]

func GET[T any](path string, params H, headers ...H) *Client[T]

func POST[T any](path string, params H, data any, headers ...H) *Client[T]

func PUT[T any](path string, params H, data any, headers ...H) *Client[T]

func DELETE[T any](path string, params H, data any, headers ...H) *Client[T]

func PATCH[T any](path string, params H, data any, headers ...H) *Client[T]

func HEAD[T any](path string, params H, headers ...H) *Client[T]

func OPTIONS[T any](path string, headers ...H) *Client[T]
```

#### JavaScript Axios-style

Example code: [try-catch-style.go](./examples/try-catch-style/main.go) &nbsp; | &nbsp; Recommended 🌟🌟🌟🌟

```textmate
gloria.GET[T]().
    Then(func(data T) {}).
    Catch(func(e *gloria.Exception) {}).
    Finally(func(c *gloria.Client[Result]) {}, bool)
```

### <span id="MoreAPI">More API Function Signatures</span>

#### Construction Function Related

```textmate
func New[T any]() *Client[T]
func Default[T any]() *Client[T]

func NewREST[T any]() *Client[T]
func NewHTTP[T any]() *Client[T]
```

#### Option Configuration Related

```textmate
func (c *Client[T]) Optional(fns ...ClientFunc[T]) *Client[T]

func Lambda[T any](f func(*Client[T])) ClientFunc[T]

func WithTimeout[T any](timeout time.Duration) ClientFunc[T]
func WithSkipTLS[T any](skipTLS bool) ClientFunc[T]
func WithFilterSlash[T any](filterSlash bool) ClientFunc[T]
func WithIsDebug[T any](isDebug bool) ClientFunc[T]
func WithUseLogger[T any](enabled bool) ClientFunc[T]
func WithModifySuccessCode[T any](code int) ClientFunc[T]

func (c *Client[T]) ToggleMode() *Client[T]                     // Toggle to the other mode.
func (c *Client[T]) FilterUrlSlash() *Client[T]                 // Trailing slashes in URLs will be automatically filtered out.
func (c *Client[T]) DefineOkCode(code int) *Client[T]           // Set a custom success value to be used as a basis for automatically determining business failures.
```

#### Request Configuration Related

```textmate
func (c *Client[T]) SetMethod(method string) *Client[T]

func (c *Client[T]) SetSchema(scheme string) *Client[T]
func (c *Client[T]) SetHost(host string) *Client[T]
func (c *Client[T]) SetBaseURI(baseUri string) *Client[T]
func (c *Client[T]) SetEndpoint(endpoint string) *Client[T]

func (c *Client[T]) SetURL(scheme, host, baseUri, endpoint string) *Client[T]

func (c *Client[T]) SetRequest(method, path string, pathParams ...string) *Client[T]

func (c *Client[T]) SetQueryParam(key, value string) *Client[T]
func (c *Client[T]) SetQueryParams(params H) *Client[T]

func (c *Client[T]) SetHeader(key, value string) *Client[T]
func (c *Client[T]) SetHeaders(headers H) *Client[T]

func (c *Client[T]) SetCookie(cookie *http.Cookie) *Client[T]
func (c *Client[T]) SetCookies(cookies []*http.Cookie) *Client[T]

func (c *Client[T]) SetBasicAuth(username, password string) *Client[T]
func (c *Client[T]) SetBearerAuth(token string) *Client[T]

func (c *Client[T]) SetAccept(accept string) *Client[T]
func (c *Client[T]) SetContentType(ct string) *Client[T]
func (c *Client[T]) SetLanguage(lang string) *Client[T]
func (c *Client[T]) SetUserAgent(ua string) *Client[T]
```

#### Middleware hooks Related

```textmate
func (c *Client[T]) UsePreHooks(funcs ...beforeRequest[T])
func (c *Client[T]) UsePostHooks(funcs ...afterResponse[T])
```

#### Request and Response handling Related

```textmate
func (c *Client[T]) Send() *Client[T]

func (c *Client[T]) Unwrap() (*Client[T], string)

func GET[T any](path string, params H, headers ...H) *Client[T]
func POST[T any](path string, params H, data any, headers ...H) *Client[T]
func PUT[T any](path string, params H, data any, headers ...H) *Client[T]
func DELETE[T any](path string, params H, data any, headers ...H) *Client[T]
func PATCH[T any](path string, params H, data any, headers ...H) *Client[T]
func HEAD[T any](path string, params H, headers ...H) *Client[T]
func OPTIONS[T any](path string, headers ...H) *Client[T]

func Request[T any](path string, params H, data any, headers ...H) ExecMethod[T]

func (c *Client[T]) Then(cb CallbackOk[T]) *Client[T]
func (c *Client[T]) Catch(cb CallbackErr) *Client[T]
func (c *Client[T]) Finally(cb CallbackExtra[T], printLog ...bool)

func (c *Client[T]) Data() T
```

#### Retrieving settings Related

```textmate
func (c *Client[T]) GetQuery(q string) string
func (c *Client[T]) GetQueryParams() SMap

func (c *Client[T]) GetHeader(key string) string
func (c *Client[T]) GetHeaders() http.Header

func (c *Client[T]) GetCookie(name string) (*http.Cookie, error)
func (c *Client[T]) GetCookies() []*http.Cookie
```

#### Print information Related

```textmate
func (c *Client[T]) Echo()

func (c *Client[T]) EchoURL() (string, string)
func (c *Client[T]) EchoCode() (int, int)
func (c *Client[T]) EchoMessage() (string, string)
func (c *Client[T]) EchoProto() string
func (c *Client[T]) EchoMode() string
func (c *Client[T]) EchoQPS() float64
func (c *Client[T]) EchoBenchmark() (int, int64)
func (c *Client[T]) EchoTime() (time.Duration, time.Time)
```

#### Logging information Related

```textmate
func (l level) ANSIColorCode() string

func (c *Client[T]) ChalkObj(level level, obj any) *Client[T]
func (c *Client[T]) ChalkStr(level level, s string) *Client[T]
func (c *Client[T]) ChalkInt(level level, n int) *Client[T]
func (c *Client[T]) ChalkPrintf(level level, format string, args ...any) *Client[T]
```

### <span id="HttpInterceptor">Injecting Interceptors for HTTP</span>

#### Injection timing

Please note the timing of loading request and response interceptors (middlewares):

```go
client := New[any]()

client.Setxxx().Setxxx()

// Add multiple request hooks
client.UsePreHooks(func(c *Client[[]ImageSearch]) error {
    fmt.Println("1. I am a request hook 🪝🪝🪝")
    return nil
}, func(c *Client[[]ImageSearch]) error {
    fmt.Println("2. I am a request hook 🐒🐒🐒")
    return nil
})

// Add multiple response hooks
client.UsePostHooks(func(c *Client[[]ImageSearch]) error {
    fmt.Println("3. I am a response hook 🪝🪝🪝")
    return nil
}, func(c *Client[[]ImageSearch]) error {
    fmt.Println("4. I am a response hook 🍌🍌🍌")
    return nil
})

// Note that the hooks should be added before calling Send!
client.Send().Unwrap()
```

#### Execution order

First, the request hooks functions will be executed in the order they were added. 
Then the request will be sent. Finally, the response hooks functions will be executed 
in the order they were added.

```textmate
1. I am a request hook 🪝🪝🪝
2. I am a request hook 🐒🐒🐒
Sending request...
3. I am a response hook 🪝🪝🪝
4. I am a response hook 🍌🍌🍌
```

### <span id="CodeSuggestions">Some Coding Suggestions</span>

#### Type Aliases

It is recommended to use predefined type aliases to maintain consistent code style.

```go
type H = map[string]any

type SMap = map[string]string

// "http://example.org/?uid=47200957&username=Mystic&is_output=true&mime_types=png,gif,ico"
// When setting request parameters, it is advisable to use the actual data types instead of relying solely on strings.

/* Good example */
c.SetQueryParams(gloria.H{
    "uid": 47200957,
    "username": "Mystic",
    "is_output": true,
    "mime_types": ["png", "gif", "ico"],
})

/* Bad example */
c.SetQueryParams(map[string]interface{}{
    "uid": "47200957",
    "username": "Mystic",
    "is_output": "true",
    "mime_types": "png,gif,ico",
})
```

#### Constants

Predefined constants are preferred over hardcoding numbers or strings directly.

```go
/* Good example */
// Use predefined constants from the gloria package, such as:
c.DefineOkCode(gloria.OkCode)

// Use predefined constants from the http package, such as:
c.DefineOkCode(http.StatusOK)

/* Bad example */
c.DefineOkCode(20000)
```


## Contribution

I warmly welcome your contribution! If you come across any areas for improvement or any issues that you would like to fix, please don't hesitate to send a pull request. I appreciate pull requests that include test cases for bug fixes or enhancements. I have put in my best effort to ensure decent code coverage, so feel free to write tests.

By the way, I am curious to hear your thoughts on `Gloria`. Please feel free to open an issue or send me an email. Your feedback means a great deal to me.


## Creator

[Pokeya Boa](https://github.com/pokeyaro)&nbsp;(<a href="mailto:pokeya.mystic@gmail.com">pokeya.mystic@gmail.com</a>)


## License

Gloria released under MIT license, refer [LICENSE](./LICENSE) file.
