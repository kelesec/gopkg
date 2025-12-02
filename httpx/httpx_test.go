package httpx

import (
	"crypto/tls"
	"github.com/valyala/fasthttp"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	client := NewClient().SetReadTimeout(10 * time.Second).
		SetWriteTimeout(10 * time.Second).
		SetMaxIdleConnDuration(10 * time.Second).
		SetMaxConnWaitTimeout(10 * time.Second).
		SetNoDefaultUserAgentHeader(false).
		SetDisablePathNormalizing(true).
		SetDisableHeaderNamesNormalizing(true).
		SetReadBufferSize(4 * 1024 * 1024).
		SetWriteBufferSize(4 * 1024 * 1024).
		SetMaxResponseBodySize(4 * 1024 * 1024).
		SetMaxConnsPerHost(256).
		SetTLSConfig(&tls.Config{
			InsecureSkipVerify: true,
		}).
		SetDial((&fasthttp.TCPDialer{
			Concurrency:      4096,
			DNSCacheDuration: time.Hour,
		}).Dial)

	resp, err := client.R().Get("https://baidu.com")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp.Status(), resp.respSize)
}

func TestFastHttpClient(t *testing.T) {
	readTimeout, _ := time.ParseDuration("500ms")
	writeTimeout, _ := time.ParseDuration("500ms")
	maxIdleConnDuration, _ := time.ParseDuration("1m")

	fc := fasthttp.Client{
		ReadTimeout:                   readTimeout,
		WriteTimeout:                  writeTimeout,
		MaxIdleConnDuration:           maxIdleConnDuration,
		NoDefaultUserAgentHeader:      true,
		DisableHeaderNamesNormalizing: true,
		DisablePathNormalizing:        true,
		ReadBufferSize:                4 * 1024 * 1024,
		WriteBufferSize:               4 * 1024 * 1024,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		Dial: (&fasthttp.TCPDialer{
			Concurrency:      4096,
			DNSCacheDuration: time.Hour,
		}).Dial,
	}

	client := NewClientWithFastHttp(&fc)
	resp, err := client.R().Get("https://baidu.com")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp.Status(), resp.respSize)
}

func TestRequest(t *testing.T) {
	req := NewClient().R()
	resp, err := req.SetMethod(MethodGet).
		SetHeaders(map[string]string{
			"header_2": "value_2",
			"header_3": "value_3",
		}).
		SetHeader("header_1", "value_1").
		SetCookies(map[string]string{
			"cookie_2": "value_2",
			"cookie_3": "value_3",
		}).
		SetCookie("cookie_1", "value_1").
		SetContentType(MIMEApplicationJSON).
		SetContentLength(2048).
		SetUserAgent(DefaultUserAgent).
		SetQueryParams(map[string]string{
			"query_2": "value_2",
			"query_3": "value_3",
		}).SetQueryParam("query_1", "value_1").
		//SetFormDatas(map[string]string{
		//	"form_2": "value_2",
		//	"form_3": "value_3",
		//}).SetFormData("form_1", "value_1").
		SetBodyString(`{"username": "admin", "password": "123123"}`).
		SetBasicAuth("username", "password").
		Do("https://baidu.com", "")

	if err != nil {
		t.Fatal(err)
	}
	t.Log(req)
	t.Log(resp)
}

func TestProxy(t *testing.T) {
	req := NewClient().SetProxy("http://8.147.118.237:18920").R()
	//req := NewClient().SetProxies([]string{
	//	"socks5://127.0.0.1:8083",
	//	"http://127.0.0.1:8083",
	//}).R()
	resp, err := req.Get("https://www.baidu.com")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp.Status(), resp.respSize)
	t.Log(resp.BodyString())
}

func TestRedirect(t *testing.T) {
	resp, err := NewClient().R().
		SetUserAgent(ChromeUserAgent).
		AllowRedirect().
		AllowSaveResponseHistory().
		Get("http://127.0.0.1:8081")

	if err != nil {
		t.Fatal(err)
	}

	for _, response := range resp.ResponseHistory() {
		t.Log(response)
	}
}
