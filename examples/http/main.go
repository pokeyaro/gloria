package main

import (
	"context"
	"fmt"

	"github.com/pokeyaro/gloria/v2"
)

type HttpBin struct {
	Slideshow struct {
		Title string `json:"title"`
	} `json:"slideshow"`
}

func main() {
	cli := gloria.NewClient[HttpBin]("https://httpbin.org")
	cli.SetRequest(gloria.MethodGet, "/json")

	if _, err := cli.SendCtx(context.Background()); err != nil {
		panic(err)
	}
	if _, err := cli.Decode(); err != nil {
		panic(err)
	}

	fmt.Println("title:", cli.Data().Slideshow.Title) // title: Sample Slide Show

	cli.Echo() // [API] GET https://httpbin.org/json -> 200 (xxx ms)
}
