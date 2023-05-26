<div align=center>
  <img src="logo.png" width="450" height="225" alt="gloria" />

  <br/>

  ![Go version](https://img.shields.io/badge/go-%3E%3Dv1.18-9cf)
  ![Release](https://img.shields.io/badge/release-1.0.0-green.svg)
  [![GoDoc](https://godoc.org/github.com/pokeyaro/gloria?status.svg)](https://godoc.org/github.com/pokeyaro/gloria)
  [![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
</div>

<p style="font-size: 15px">
  Gloria 是一款旨在简洁、优雅，主打 RESTful 风格的 HTTP 客户端工具。Gloria 启示 Go 网络请求之美！
  灵感来源：Gloria 受到了 <a href="https://github.com/go-resty">Go resty</a>，<a href="https://github.com/requests">Python requests</a> 和 <a href="https://github.com/axios">JavaScript axios</a> 的启发。
</p>

简体中文 | [English](./README.md)


## 由来

`"Gloria"` 一词取自歌手邓紫棋专辑《启示录》的[主打歌](https://www.youtube.com/watch?v=stGUpMav1sc)，我们希望这个库能像邓紫棋的音乐一样，给人们带来美妙的体验。

同时，`Gloria` 在拉丁语中意为“荣耀”，寓意着这个库的设计旨在提供一种荣耀的、优雅的方式来使用 `Go` 的 `HTTP` 客户端，就像 `RESTful` 规范一样，使得网络请求更加符合人类的自然习惯。

我们向所有为构建互联网应用程序而努力工作的人们致敬，希望 `Gloria` 能够为他们带来更加舒适的开发及使用体验。


## 特性

- 🪶 简洁、易用、丰富多样的 `API` 设计
- 👏 泛型实现统一 `RESTful` 风格响应
- 🚀 支持多个请求、响应拦截器的注入
- 📝 详细的彩印日志以及错误追踪定位功能
- 💃🏼 优雅的 `Unwrap` 错误处理（`Rust` 风格）
- 🧭 请求调用时间及 `QPS` 预估功能
- 🌈 `GET`、`POST`、`PUT`、`DELETE` 等
- 🎋 经过良好的 `Benchmark` 测试等


## 用法

### 版本

支持的版本：依赖 `Go` `T` 泛型特性，至少要求解释器 `go1.18+` 以上。

### 安装

```bash
go get -u github.com/pokeyaro/gloria
```

### 导入

```go
import "github.com/pokeyaro/gloria"
```


## 索引

- [简单的 GET 请求](#SimpleGet)
- [更丰富的请求配置](#ExtendedConf)
- [常用 CRUD 请求操作](#CrudMethod)
- [两种响应模式](#TwoRespMode)
- [REST 语法糖](#RestSyntacticSugar)
- [更多 API 函数签名](#MoreAPI)
- [为 HTTP 注入拦截器](#HttpInterceptor)
- [一些代码编写建议](#CodeSuggestions)


## 文档

通过下面的示例，能够快速入门了解 `gloria` 库

### <span id="SimpleGet">简单的 GET 请求</span>

> 真实 API 可使用：`curl -X GET --location 'http://httpbin.org/get'`

#### 示例一

首先，为了更清晰的演示该如何使用 `gloria`，我们使用 `NewHTTP()` 方法，通过分步操作来进行说明：

示例代码：[api-httpbin.go](./example_test.go) &nbsp; | &nbsp; 推荐 🌟🌟🌟

```go
// 准备响应体的结构体
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

// 使用 NewHTTP 方法进行构建请求
client := gloria.NewHTTP[HttpBin]()

// 设置请求类型
client.SetMethod(gloria.MethodGet)

// 设置请求的请求路由资源（URL分段填写：proto, host, baseURI, endpoint）
client.SetURL(gloria.ProtocolHttp, "httpbin.org", "", "/get")

// 发送请求
client.Send()

// 接住错误
client.Unwrap()

// 输出请求元信息
client.Echo()

/* 输出
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

#### 示例二

为用户提供了更为简洁的 `API` 方法，并且支持链式加载，于是你可以写成下面这样：

示例代码：[api-httpbin.go](./example_test.go) &nbsp; | &nbsp; 推荐 🌟🌟🌟🌟🌟

```go
type HttpBin struct {
    // 省略...
}

// 请求
client, _ := gloria.NewHTTP[HttpBin]().SetRequest(gloria.MethodGet, "http://httpbin.org/get").Send().Unwrap()

// 获取
fmt.Println(client.Data().Url)

/* 输出
http://httpbin.org/get
*/
```

### <span id="ExtendedConf">更丰富的请求配置</span>

> 真实 API 可使用：`curl -X GET --location 'https://api.thecatapi.com/v1/images/search?size=med&mime_types=png%2Cgif&format=json&order=RANDOM&limit=20' --header 'Content-Type: application/json'`

#### 示例一

使用 `New()` 和 `Optional()` 方法来加载更丰富的配置，我们先来感受一下吧：

示例代码：[api-cat.go](./examples/request_test.go) &nbsp; | &nbsp; 推荐 🌟🌟🌟🌟

```go
// 准备响应体的结构体
type ImageSearch struct {
    Id     string `json:"id"`
    Url    string `json:"url"`
    Width  int    `json:"width"`
    Height int    `json:"height"`
}

// 构建一个更丰富的请求
func GetCatAPI() {
	// 这个接口返回的是结构体切片
    r := gloria.New[[]ImageSearch]()

    r.Optional(
        // 是否开启 debug 模式
        gloria.WithIsDebug[[]ImageSearch](false),
        // 是否启动日志打印
        gloria.WithUseLogger[[]ImageSearch](true),
        // 通过 lambda 方式设置更多参数
        gloria.Lambda[[]ImageSearch](func(c *gloria.Client[[]ImageSearch]) {
            c.Config.SkipTLS = true
            c.Config.Timeout = gloria.TimeoutMedium
            c.Config.IsRestMode = false // 如果非标准 RESTful 响应接口，使用原生 HTTP 模式，即选择 false
        }),
    ).
        // 设置请求方法
        SetMethod(gloria.MethodGet).
        // 设置请求路径，该方法需分段编写
        SetURL(gloria.ProtocolHttps, "api.thecatapi.com", "/v1", "/images/search").
        // 设置多个请求参数
		SetQueryParams(gloria.H{
            "size":       "med",
            "mime_types": []string{"png", "gif"},
            "format":     "json",
            "order":      "RANDOM",
            "limit":      20,
        }).
        // 设置多个请求头
        SetHeaders(gloria.H{
            "x-api-key": "live_example-api-key",
            "Content-Type": "application/json",
        }).
        // 发送请求
        Send().
        Unwrap()

    fmt.Println(r.Data()[0].Url)	
}
```

#### 示例二

使用 `Default()` 来加载默认配置，事实上你根本不需要关心 `Default()` 与 `New()` 的区别，或者它内部增加了哪些默认参数，如果你特别想知道，可通过 `go doc gloria.Default` 命令查看：

示例代码：[api-cat.go](./examples/request_test.go) &nbsp; | &nbsp; 推荐 🌟🌟🌟🌟🌟

```go
type ImageSearch struct {
    // 省略...
}

// 构建一个更丰富的请求
func GetCatAPI() {
	// 注意 Default 方法模式是 REST 模式，这里需要使用 ToggleMode 进行切换！
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

### <span id="CrudMethod">常用 CRUD 请求操作</span>

下面演示 `API` 的 `CRUD` 操作，分别为 `Create` `[POST]`、 `Read` `[GET]`、`Update` `[PUT]`、`Delete` `[DELETE]`

#### [GET] 请求示例

示例代码：[api-cat.go](./examples/request_test.go) &nbsp; | &nbsp; 推荐 🌟🌟🌟🌟🌟

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

#### [POST] 请求示例

示例代码：[api-cat.go](./examples/request_test.go) &nbsp; | &nbsp; 推荐 🌟🌟🌟🌟🌟

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
        SubId:   "my-key-123445",
    }

    r := gloria.NewHTTP[FavouriteImgResp]()

    r.SetRequest(gloria.MethodPost, "https://api.thecatapi.com/v1/favourites").SetHeaders(gloria.H{
        "x-api-key":    "your-api-key",
        "Content-Type": "application/json",
    }).SetPayload(&data).Send().Unwrap()

    fmt.Println("post_id:", r.Data().Id)
}
```

#### [PUT] 请求示例

暂无示例，可参考 `POST` 方法。

#### [DELETE] 请求示例

示例代码：[api-cat.go](./examples/request_test.go) &nbsp; | &nbsp; 推荐 🌟🌟🌟🌟🌟

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

### <span id="TwoRespMode">两种响应模式</span>

拥有两种不同的响应模式：`HTTP` 模式和 `REST` 模式，下面我们来分别解释它们的区别！

所谓 `HTTP` 模式，可以是任何的响应数据！就像上面例子那样，没有任何规范可言！自由度较高！

所谓 `REST` 模式，如果你的响应格式形如下方：

```json
{
  "code": 0,
  "msg": "success",
  "data": Object,
}
```

知道了吧，我们 `REST` 模式其实就是直接解析 `data` 中的数据，让你的数据请求一步到位，响应反馈的结果可以更加直接！

下面我们用一个表格来总结其 `API` 用法：

<table>
  <tr>
    <th width="15%">响应格式</th>
    <th width="35%">选择方法</th>
    <th width="auto">场景说明</th>
  </tr>
  <tr>
    <td>HTTP 模式</td>
    <td>
      <span>方式1:</span>&nbsp;<code>New().ToggleMode()</code> <br/>
      <span>方式2:</span>&nbsp;<code>Default().ToggleMode()</code> <br/>
      <span>方式3:</span>&nbsp;<code>NewHTTP()</code> // 本质: 方式1的语法糖 <br/>
    </td>
    <td>如果你的响应格式为非标准的 RESTful 响应接口，请使用更通用的 HTTP 模式。</td>
  </tr>
  <tr>
    <td>REST 模式</td>
    <td>
      <span>方式1:</span>&nbsp;<code>New()</code> <br/>
      <span>方式2:</span>&nbsp;<code>Default()</code> <br/>
      <span>方式3:</span>&nbsp;<code>NewREST()</code> // 本质: 方式1的语法糖 <br/>
    </td>
    <td>当你的响应格式为标准的 RESTful 响应接口，形如：<code>{"code": 0, "msg": "success", "data": null}</code></td>
  </tr>
</table>

### <span id="RestSyntacticSugar">REST 语法糖</span>

> 优势：写法更为简洁，所有配置均在该函数内完成解析并发送请求！因此，即是优点，也是缺点。
> 
> 缺点：如上面所说，暂不支持更高级的语法，如：这两种语法糖会直接发送请求，因此，无法链式注入其他设置！

当我们使用 **REST 模式** 来构建处理请求，提供额外的两种非常易用的语法糖！

#### 类 Python Requests 风格

示例代码：[requests-style.go](./examples/requests-style/main.go) &nbsp; | &nbsp; 推荐 🌟🌟🌟🌟🌟

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

#### 类 Javascript Axios 风格

示例代码：[try-catch-style.go](./examples/try-catch-style/main.go) &nbsp; | &nbsp; 推荐 🌟🌟🌟🌟

```textmate
gloria.GET[T]().
    Then(func(data T) {}).
    Catch(func(e *gloria.Exception) {}).
    Finally(func(c *gloria.Client[Result]) {}, bool)
```

### <span id="MoreAPI">更多 API 函数签名</span>

#### 构建函数

```textmate
func New[T any]() *Client[T]
func Default[T any]() *Client[T]

func NewREST[T any]() *Client[T]
func NewHTTP[T any]() *Client[T]
```

#### 选项配置

```textmate
func (c *Client[T]) Optional(fns ...ClientFunc[T]) *Client[T]

func Lambda[T any](f func(*Client[T])) ClientFunc[T]

func WithTimeout[T any](timeout time.Duration) ClientFunc[T]
func WithSkipTLS[T any](skipTLS bool) ClientFunc[T]
func WithFilterSlash[T any](filterSlash bool) ClientFunc[T]
func WithIsDebug[T any](isDebug bool) ClientFunc[T]
func WithUseLogger[T any](enabled bool) ClientFunc[T]
func WithModifySuccessCode[T any](code int) ClientFunc[T]

func (c *Client[T]) ToggleMode() *Client[T]                     // 切换到另外一种模式。
func (c *Client[T]) FilterUrlSlash() *Client[T]                 // 将自动过滤掉URL尾部的斜杠。
func (c *Client[T]) DefineOkCode(code int) *Client[T]           // 设置自定义成功返回值，作为用于自动判断业务失败的依据。
```

#### 请求设置

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

#### 中间件钩子

```textmate
func (c *Client[T]) UsePreHooks(funcs ...beforeRequest[T])
func (c *Client[T]) UsePostHooks(funcs ...afterResponse[T])
```

#### 请求响应

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

#### 获取设置

```textmate
func (c *Client[T]) GetQuery(q string) string
func (c *Client[T]) GetQueryParams() SMap

func (c *Client[T]) GetHeader(key string) string
func (c *Client[T]) GetHeaders() http.Header

func (c *Client[T]) GetCookie(name string) (*http.Cookie, error)
func (c *Client[T]) GetCookies() []*http.Cookie
```

#### 输出信息

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

#### 日志

```textmate
func (l level) ANSIColorCode() string

func (c *Client[T]) ChalkObj(level level, obj any) *Client[T]
func (c *Client[T]) ChalkStr(level level, s string) *Client[T]
func (c *Client[T]) ChalkInt(level level, n int) *Client[T]
func (c *Client[T]) ChalkPrintf(level level, format string, args ...any) *Client[T]
```

### <span id="HttpInterceptor">为 HTTP 注入拦截器</span>

#### 注入时机

请注意加载请求、响应拦截器（中间件）的时机：

```go
client := New[any]()

client.Setxxx().Setxxx()

// 添加多个请求钩子
client.UsePreHooks(func(c *Client[[]ImageSearch]) error {
    fmt.Println("1. 我是一只请求 🪝🪝🪝")
    return nil
}, func(c *Client[[]ImageSearch]) error {
    fmt.Println("2. 我是一只请求 🐒🐒🐒")
    return nil
})

// 添加多个响应钩子
client.UsePostHooks(func(c *Client[[]ImageSearch]) error {
    fmt.Println("3. 我是一只响应 🪝🪝🪝")
    return nil
}, func(c *Client[[]ImageSearch]) error {
    fmt.Println("4. 我是一只响应 🍌🍌🍌")
    return nil
})

// 注意，需要在 Send 之前完成钩子的注入！
client.Send().Unwrap()
```

#### 执行顺序

首先会执行请求钩子函数，按照它们添加的顺序依次执行。然后发送请求。最后执行响应钩子函数，同样按照添加的顺序依次执行。

```textmate
1. 我是一只请求 🪝🪝🪝
2. 我是一只请求 🐒🐒🐒
发送请求...
3. 我是一只响应 🪝🪝🪝
4. 我是一只响应 🍌🍌🍌
```

### <span id="CodeSuggestions">一些代码编写建议</span> 

#### 类型别名

更推荐使用预定义类型别名来统一代码风格！

```go
type H = map[string]any

type SMap = map[string]string

// "http://example.org/?uid=47200957&username=Mystic&is_output=true&mime_types=png,gif,ico"
// 当你需要设置请求参数时，请让你的代码value类型回归本真，而不是全部使用string类型！

/* 正例 */
c.SetQueryParams(gloria.H{
    "uid": 47200957,
    "username": "Mystic",
    "is_output": true,
    "mime_types": ["png", "gif", "ico"],
})

/* 反例 */
c.SetQueryParams(map[string]interface{}{
    "uid": "47200957",
    "username": "Mystic",
    "is_output": "true",
    "mime_types": "png,gif,ico",
})
```

#### 常量

更推荐使用预定义的常量，而不是直接硬编码数字或字符串！

```go
/* 正例 */
// 使用预定义常量，如：
c.DefineOkCode(gloria.OkCode)

// 使用http库预定义常量，如：
c.DefineOkCode(http.StatusOK)

/* 反例 */
c.DefineOkCode(20000)
```


## 贡献

非常欢迎您的贡献！如果您发现任何改进或需要修复的问题，请随时提交拉取请求。我喜欢包含针对修复或增强的测试用例的拉取请求。我已经尽力提供了相当不错的代码覆盖率，所以请随意编写测试。

顺便说一下，我很想知道您对 `Gloria` 的看法。请随时提交问题或给我发送电子邮件；这对我来说非常重要。


## 作者

[Pokeya Boa](https://github.com/pokeyaro)&nbsp;(<a href="mailto:pokeya.mystic@gmail.com">pokeya.mystic@gmail.com</a>)


## 许可证

`Gloria` 使用 MIT 许可证进行发布，详见 [LICENSE](./LICENSE) 文件。
