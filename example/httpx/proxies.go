package main

import (
	"fmt"
	"github.com/kelesec/gopkg/httpx"
	"log"
)

func main() {
	// 设置代理
	req := httpx.NewClient().SetProxies([]string{
		"http://127.0.0.1:8083",
	}).SetProxy("socks5://127.0.0.1:8083").R()

	// 允许重定向请求，并记录历史响应
	resp, err := req.SetUserAgent(httpx.ChromeUserAgent).
		SetBasicAuth("zhangsan", "123456").
		AllowRedirect().
		AllowSaveResponseHistory().
		Get("http://127.0.0.1:8081")
	if err != nil {
		log.Fatal(err)
	}

	for _, rsp := range resp.ResponseHistory() {
		fmt.Println(rsp)
	}
}
