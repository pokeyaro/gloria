// Copyright (c) 2023 Pokeya Boa <pokeya.mystic@gmail.com>, All rights reserved.
// resty source code and usage is governed by a MIT style
// license that can be found in the LICENSE file.

package gloria

import (
	"fmt"
)

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

func ExampleNewHTTP_method1() {
	// Create an HTTP client with HttpBin response struct
	client := NewHTTP[HttpBin]()

	client.SetMethod(MethodGet)
	client.SetURL(ProtocolHttp, "httpbin.org", "", "/get")
	client.Send()
	client.Unwrap()

	client.Echo()
}

func ExampleNewHTTP_method2() {
	// One line response
	client, _ := NewHTTP[HttpBin]().SetRequest(MethodGet, "http://httpbin.org/get").Send().Unwrap()

	fmt.Println(client.Data().Url)

	// Output:
	// http://httpbin.org/get
}
