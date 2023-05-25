// Copyright (c) 2023 Pokeya Boa <pokeya.mystic@gmail.com>, All rights reserved.
// resty source code and usage is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"time"

	"github.com/pokeyaro/gloria"
)

/*
The parsed response body is as follows (mock data):
{
    "code":0,
    "msg":"success",
    "data":{
        "items":[
            {
                "id":13473,
                "ip_addr":"10.98.100.132",
                "status":"ok",
                "version":"v9.0.0",
                "site_name":"1234 Main Street, Anytown, USA",
                "ctime":"2022-01-27T18:36:02+08:00",
                "utime":"2023-05-20T18:00:10+08:00",
            },
            {
                "id":13474,
                "ip_addr":"10.98.100.133",
                "status":"error",
                "version":"v9.0.0",
                "site_name":"5678 Elm Avenue, Somewhere City, Canada",
                "ctime":"2022-01-27T18:36:02+08:00",
                "utime":"2023-05-20T18:00:10+08:00",
            }
        ],
        "page_info":{
            "page_num":1,
            "page_size":2,
            "total_count":1086
        }
    }
}
*/

type Items struct {
	ID         int       `json:"id"`
	IpAddr     string    `json:"ip_addr"`
	Status     string    `json:"status"`
	Version    string    `json:"version"`
	SiteName   string    `json:"site_name"`
	CreateTime time.Time `json:"ctime"`
	UpdateTime time.Time `json:"utime"`
}

type PageInfo struct {
	PageNum    int `json:"page_num"`
	PageSize   int `json:"page_size"`
	TotalCount int `json:"total_count"`
}

type Result struct {
	Items    []Items  `json:"items"`
	PageInfo PageInfo `json:"page_info"`
}

func main() {
	// Need to add an extra layer!
	type JsonResponse struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data Result `json:"data"`
	}

	// If HTTP mode is used, the response format is the original response body!
	r := gloria.NewHTTP[JsonResponse]()
	r.
		SetRequest(gloria.MethodGet, "https://example.org/api/v1/meta/address?page_num=1&page_size=2").
		SetHeaders(gloria.H{
			"Uid":      47200957,
			"Username": "mystic",
			"X-Token":  "3de4fd33-f18b-4e2a-906b-7927f6e5828f",
		}).
		Send().Unwrap()

	resp := r.Data()

	fmt.Println(resp.Code)
	fmt.Println(resp.Msg)
	fmt.Println(resp.Data)
}
