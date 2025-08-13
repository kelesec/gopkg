package httpx

import (
	"github.com/valyala/fasthttp"
)

type Header map[string]string

func (h Header) Get(key string) string {
	return h[key]
}

type Response struct {
	OriginalRequest  fasthttp.Request
	OriginalResponse fasthttp.Response

	header          Header      // 响应头
	headerBytes     []byte      // 响应头字节
	body            []byte      // 响应体
	contentLength   int         // 响应体长度
	respSize        int         // 响应长度（响应头+响应体）
	location        string      // 30X跳转后的地址
	responseHistory []*Response // 允许重定向跳转时，记录每次请求的响应，包括最后一次请求也会记录
}

func (r *Response) Status() int {
	return r.OriginalResponse.StatusCode()
}

func (r *Response) Header() Header {
	return r.header
}

func (r *Response) HeaderBytes() []byte {
	return r.headerBytes
}

func (r *Response) HeaderString() string {
	return string(r.headerBytes)
}

func (r *Response) Body() []byte {
	return r.body
}

func (r *Response) BodyString() string {
	return string(r.body)
}

func (r *Response) ContentLength() int {
	return r.contentLength
}

func (r *Response) ResponseSize() int {
	return r.respSize
}

func (r *Response) Location() string {
	return r.location
}

func (r *Response) ResponseHistory() []*Response {
	return r.responseHistory
}

func (r *Response) String() string {
	return r.OriginalResponse.String()
}
