// Copyright (c) 2023 Pokeya Boa <pokeya.mystic@gmail.com>, All rights reserved.
// resty source code and usage is governed by a MIT style
// license that can be found in the LICENSE file.

package examples

import (
	"fmt"
	"testing"

	"gloria"
)

type ImageSearch struct {
	Id     string `json:"id"`
	Url    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

/*
	curl -X GET --location 'https://api.thecatapi.com/v1/images/search?order=RANDOM&limit=20&size=med&mime_types=png,gif&format=json' \
	--header 'Content-Type: application/json' \
	--header 'x-api-key: live_example-api-key'
*/

func TestRequest_New(t *testing.T) {
	r := gloria.New[[]ImageSearch]()

	r.Optional(
		// single setting
		gloria.WithIsDebug[[]ImageSearch](false),
		gloria.WithUseLogger[[]ImageSearch](true),
		gloria.Lambda[[]ImageSearch](func(c *gloria.Client[[]ImageSearch]) {
			// More settings
			c.Config.SkipTLS = true
			c.Config.Timeout = gloria.TimeoutMedium
			c.Config.IsRestMode = false

		}),
	).
		SetMethod(gloria.MethodGet).
		SetURL(gloria.ProtocolHttps, "api.thecatapi.com", "/v1", "/images/search").
		SetQueryParams(gloria.H{
			"size":       "med",
			"mime_types": []string{"png", "gif"},
			"format":     "json",
			"order":      "RANDOM",
			"limit":      20,
		}).
		SetHeaders(gloria.H{
			"x-api-key":    "live_example-api-key",
			"Content-Type": "application/json",
		}).
		Send().Unwrap()

	fmt.Println(r.Result.Data[0].Url)
}

func TestRequest_Default(t *testing.T) {
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
		SetLanguage(gloria.LocaleEn).
		Send().Unwrap()

	fmt.Println(r.Result.Data[0].Url)
}
