// Copyright (c) 2023 Pokeya Boa <pokeya.mystic@gmail.com>, All rights reserved.
// resty source code and usage is governed by a MIT style
// license that can be found in the LICENSE file.

package examples

import (
	"fmt"
	"testing"
	"time"

	"github.com/pokeyaro/gloria"
)

/*
	curl -X GET --location 'https://api.thecatapi.com/v1/images/search?order=RANDOM&limit=20&size=med&mime_types=png,gif&format=json' \
	--header 'Content-Type: application/json' \
	--header 'x-api-key: your-api-key'
*/

type ImageSearch struct {
	Id     string `json:"id"`
	Url    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

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
			c.Config.FilterSlash = true
			c.Config.IsRestMode = false

		}),
	).
		SetMethod(gloria.MethodGet).
		SetURL(gloria.ProtocolHttps, "api.thecatapi.com", "/v1", "/images/search/").
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

/*
	curl -X GET --location 'https://api.thecatapi.com/v1/favourites' \
	--header 'Content-Type: application/json' \
	--header 'x-api-key: your-api-key'
*/

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

func TestRequest_Get(t *testing.T) {
	r := gloria.NewHTTP[[]FavouritesList]()

	r.SetRequest(gloria.MethodGet, "https://api.thecatapi.com/v1/favourites").SetHeaders(gloria.H{
		"x-api-key":    "your-api-key",
		"Content-Type": "application/json",
	}).Send().Unwrap()

	for _, v := range r.Data() {
		fmt.Println(v)
	}
}

/*
	curl -X POST --location 'https://api.thecatapi.com/v1/favourites' \
	--header 'x-api-key: your-api-key' \
	--header 'Content-Type: application/json' \
	--data '{
	  "image_id": "12345",
	  "sub_id": "my-key-12345"
	}'
*/

type FavouriteImgResp struct {
	Message string `json:"message"`
	Id      int    `json:"id"`
}

type FavouriteImgBody struct {
	ImageId string `json:"image_id"`
	SubId   string `json:"sub_id"`
}

func TestRequest_Post(t *testing.T) {
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

/*
	curl -X DELETE --location 'https://api.thecatapi.com/v1/favourites/232338740' \
	--header 'Content-Type: application/json' \
	--header 'x-api-key: your-api-key'
*/

type Result struct {
	Message string `json:"message"`
}

func TestRequest_Delete(t *testing.T) {
	r := gloria.NewHTTP[Result]()

	r.SetRequest(gloria.MethodDelete, "https://api.thecatapi.com/v1/favourites/:id", "232338734").SetHeaders(gloria.H{
		"x-api-key":    "your-api-key",
		"Content-Type": "application/json",
	}).Send().Unwrap()

	fmt.Println("message:", r.Data().Message)
}
