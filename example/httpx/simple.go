package main

import (
	"fmt"
	"github.com/kelesec/gopkg/httpx"
	"time"
)

func main() {
	client := httpx.NewClient()
	req := client.SetWriteTimeout(2 * time.Second).R()
	resp, err := req.SetQueryParam("rsv_iqid", "0xd3b3b109000a0c20").
		SetHeaders(map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		}).
		SetUserAgent(httpx.DefaultUserAgent).
		SetQueryParams(map[string]string{
			"issp":    "1",
			"rsv_bp":  "1",
			"rsv_idx": "2",
		}).Get("https://www.baidu.com/s?wd=python&rsv_spt=1")
	if err != nil {
		panic(err)
	}

	fmt.Println(req)
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println(resp)
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println(resp.Status())
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println(resp.HeaderString())
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println(resp.BodyString())
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println(resp.ContentLength())
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println(resp.ResponseSize())
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println(resp.Location())
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println(resp.Header())
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println(resp.Header().Get("Content-Type"))
}
