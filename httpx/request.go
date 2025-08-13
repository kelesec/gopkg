package httpx

import (
	"bytes"
	"fmt"
	"github.com/valyala/fasthttp"
	_url "net/url"
	"sync"
)

type Request struct {
	Method          string                 // 请求方法
	url             string                 // 请求url
	schema          string                 // 请求协议
	hostPort        string                 // 请求主机/域名+端口
	hostname        string                 // 请求主机/域名
	port            string                 // 请求端口
	path            string                 // 请求路径
	Headers         fasthttp.RequestHeader // 请求头
	Cookies         map[string]string      // Cookie
	ContentType     string                 // 请求体类型
	ContentLength   int                    // 请求体长度，默认会自动计算
	UserAgent       string                 // 请求UA
	QueryParam      _url.Values            // URL 请求参数
	FormData        _url.Values            // form-data 请求体
	Body            []byte                 // 请求体
	BasicAuth       *BasicAuth             // basic 基础认证
	OriginalRequest fasthttp.Request       // 原始请求的数据备份
	client          *Client

	allowRedirect            bool // 设置允许重定向、默认 false，也即不进行重定向请求
	allowSaveResponseHistory bool // 不保存重定向请求历史响应信息
	maxRedirectsCount        int  // 最大重定向请求次数，默认为5

	// 加个锁
	clock *sync.Mutex
}

func newRequest(client *Client) *Request {
	return &Request{
		client:                   client,
		clock:                    &sync.Mutex{},
		allowRedirect:            false,
		allowSaveResponseHistory: false,
		maxRedirectsCount:        5,
	}
}

// parseUrl 解析URL
func (r *Request) parseUrl(url string) error {
	u, err := _url.Parse(url)
	if err != nil {
		return fmt.Errorf("parse uri failed: %s", err)
	}

	for k, vs := range u.Query() {
		for _, v := range vs {
			r.QueryParam.Add(k, v)
		}
	}
	u.RawQuery = r.QueryParam.Encode()

	r.schema = u.Scheme
	r.hostPort = u.Host
	r.hostname = u.Hostname()
	r.port = u.Port()
	r.path = u.Path
	r.url = u.String()
	return nil
}

// preCheck 前置检查，主要用于将 `Request` 属性值同步给 `fasthttp.Request`
func (r *Request) preCheck(url, method string, req *fasthttp.Request) error {
	r.clock.Lock()
	defer r.clock.Unlock()

	if url == "" {
		return fmt.Errorf("URI is empty")
	} else if err := r.parseUrl(url); err != nil {
		return fmt.Errorf("parse URI %s error: %v", r.url, err)
	}

	r.Headers.CopyTo(&req.Header)
	req.SetRequestURI(r.url)
	req.Header.SetMethod(method)
	r.Method = method

	for k, v := range r.Cookies {
		req.Header.SetCookie(k, v)
	}

	if r.ContentType != "" {
		req.Header.SetContentType(r.ContentType)
	}

	if r.FormData != nil || len(r.Body) != 0 {
		if r.FormData != nil {
			r.ContentLength = len(r.FormData.Encode())
			req.Header.SetContentLength(r.ContentLength)
			req.SetBodyString(r.FormData.Encode())
		} else {
			r.ContentLength = len(r.Body)
			req.Header.SetContentLength(r.ContentLength)
			req.SetBody(r.Body)
		}
	}

	if r.UserAgent != "" {
		req.Header.SetUserAgent(r.UserAgent)
	}

	if r.BasicAuth != nil {
		header, auth := r.BasicAuth.GetBasicAuth()
		req.Header.Set(header, auth)
	}

	// 做个保存备份
	req.CopyTo(&r.OriginalRequest)
	return nil
}

// postCheck 后置检查，主要用于将 `fasthttp.Response` 属性同步给自定义的 Response
func (r *Request) postCheck(resp *fasthttp.Response) *Response {
	newResp := &Response{}
	resp.CopyTo(&newResp.OriginalResponse)
	r.OriginalRequest.CopyTo(&newResp.OriginalRequest)

	newResp.headerBytes = resp.Header.Header()
	newResp.body = resp.Body()
	newResp.respSize = len(resp.Header.Header()) + len(resp.Body())
	newResp.contentLength = resp.Header.ContentLength()
	if newResp.contentLength < 0 {
		newResp.contentLength = len(newResp.body)
	}

	if location := resp.Header.Peek("location"); len(location) > 0 {
		newResp.location = string(location)
	} else if location = resp.Header.Peek("Location"); len(location) > 0 {
		newResp.location = string(location)
	}

	newResp.header = make(Header)
	headerLines := bytes.Split(newResp.headerBytes, []byte("\n"))
	for _, line := range headerLines {
		tmp := bytes.SplitN(bytes.TrimSpace(line), []byte(":"), 2)
		if len(tmp) != 2 {
			continue
		}
		newResp.header[string(tmp[0])] = string(bytes.TrimSpace(tmp[1]))
	}

	return newResp
}

// Do 执行HTTP请求
func (r *Request) Do(url, method string) (*Response, error) {
	if method == "" {
		if r.Method != "" {
			method = r.Method
		} else {
			return nil, fmt.Errorf("method is empty")
		}
	}

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseResponse(resp)
		fasthttp.ReleaseRequest(req)
	}()

	if err := r.preCheck(url, method, req); err != nil {
		return nil, fmt.Errorf("preCheck err: %v", err)
	}

	redirectCount := 0
	finalResp := new(Response)
	respHistory := make([]*Response, 0)

	for {
		err := r.client.execute(req, resp)
		if err != nil {
			return nil, fmt.Errorf("get %s err: %v", r.url, err)
		}

		// 不允许重定向时直接退出
		if !r.allowRedirect {
			return r.postCheck(resp), nil
		}

		// 非重定向请求直接退出循环
		statusCode := resp.Header.StatusCode()
		if !fasthttp.StatusCodeIsRedirect(statusCode) {
			finalResp = r.postCheck(resp)
			respHistory = append(respHistory, finalResp)
			if len(respHistory) != 0 {
				finalResp.responseHistory = respHistory
			}
			break
		}

		// 超过最大重定向请求次数支持
		redirectCount++
		if redirectCount > r.maxRedirectsCount {
			return nil, fasthttp.ErrTooManyRedirects
		}

		tmpResp := r.postCheck(resp)
		if tmpResp.Location() == "" {
			return nil, fasthttp.ErrMissingLocation
		}

		// 保存历史请求
		if r.allowSaveResponseHistory {
			respHistory = append(respHistory, tmpResp)
		}

		// 继续重定向
		if string(req.Header.Method()) == "POST" && (statusCode == 301 || statusCode == 302) {
			req.Header.SetMethod(MethodGet)
		}
		req.SetRequestURI(tmpResp.location)
	}

	return finalResp, nil
}

func (r *Request) Get(url string) (*Response, error) {
	return r.Do(url, MethodGet)
}

func (r *Request) Head(url string) (*Response, error) {
	return r.Do(url, MethodHead)
}

func (r *Request) Post(url string) (*Response, error) {
	return r.Do(url, MethodPost)
}

func (r *Request) Put(url string) (*Response, error) {
	return r.Do(url, MethodPut)
}

func (r *Request) Patch(url string) (*Response, error) {
	return r.Do(url, MethodPatch)
}

func (r *Request) Delete(url string) (*Response, error) {
	return r.Do(url, MethodDelete)
}

func (r *Request) Connect(url string) (*Response, error) {
	return r.Do(url, MethodConnect)
}

func (r *Request) Options(url string) (*Response, error) {
	return r.Do(url, MethodOptions)
}

func (r *Request) Trace(url string) (*Response, error) {
	return r.Do(url, MethodTrace)
}

// SetMethod 设置请求方法，适合自定义请求方法的情况，需要通过 `Do` 方法发送才能生效
func (r *Request) SetMethod(method string) *Request {
	r.clock.Lock()
	defer r.clock.Unlock()
	r.Method = method
	return r
}

// Url 获取请求URL
func (r *Request) Url() string {
	r.clock.Lock()
	defer r.clock.Unlock()
	return r.url
}

// Schema 获取请求协议
func (r *Request) Schema() string {
	r.clock.Lock()
	defer r.clock.Unlock()
	return r.schema
}

// HostPort 获取目标主机/域名+端口
func (r *Request) HostPort() string {
	r.clock.Lock()
	defer r.clock.Unlock()
	return r.hostPort
}

// Hostname 获取目标主机/域名
func (r *Request) Hostname() string {
	r.clock.Lock()
	defer r.clock.Unlock()
	return r.hostname
}

// Port 获取目标端口
func (r *Request) Port() string {
	r.clock.Lock()
	defer r.clock.Unlock()
	return r.port
}

// Path 获取请求路径
func (r *Request) Path() string {
	r.clock.Lock()
	defer r.clock.Unlock()
	return r.path
}

// SetHeader 设置请求头
func (r *Request) SetHeader(key, value string) *Request {
	r.clock.Lock()
	defer r.clock.Unlock()
	r.Headers.Set(key, value)
	return r
}

// SetHeaders 设置请求头
func (r *Request) SetHeaders(headers map[string]string) *Request {
	for k, v := range headers {
		r.SetHeader(k, v)
	}
	return r
}

// SetCookie 设置 Cookie
func (r *Request) SetCookie(key, value string) *Request {
	r.clock.Lock()
	defer r.clock.Unlock()
	if r.Cookies == nil {
		r.Cookies = make(map[string]string)
	}
	r.Cookies[key] = value
	return r
}

// SetCookies 设置 Cookie
func (r *Request) SetCookies(cookies map[string]string) *Request {
	for k, v := range cookies {
		r.SetCookie(k, v)
	}
	return r
}

// SetContentType 设置请求体类型
func (r *Request) SetContentType(contentType string) *Request {
	r.clock.Lock()
	defer r.clock.Unlock()
	r.ContentType = contentType
	return r
}

// SetContentLength 设置请求体长度
func (r *Request) SetContentLength(contentLength int) *Request {
	r.clock.Lock()
	defer r.clock.Unlock()
	r.ContentLength = contentLength
	return r
}

// SetUserAgent 设置请求UA
func (r *Request) SetUserAgent(ua string) *Request {
	r.clock.Lock()
	defer r.clock.Unlock()
	r.UserAgent = ua
	return r
}

// SetQueryParam 设置URL请求参数
func (r *Request) SetQueryParam(key, value string) *Request {
	r.clock.Lock()
	defer r.clock.Unlock()
	if r.QueryParam == nil {
		r.QueryParam = _url.Values{}
	}
	r.QueryParam.Set(key, value)
	return r
}

// SetQueryParams 设置URL请求参数
func (r *Request) SetQueryParams(params map[string]string) *Request {
	for k, v := range params {
		r.SetQueryParam(k, v)
	}
	return r
}

// SetFormData 设置 form-data 请求体参数
func (r *Request) SetFormData(key, value string) *Request {
	r.clock.Lock()
	defer r.clock.Unlock()
	if r.FormData == nil {
		r.FormData = _url.Values{}
	}
	r.FormData.Set(key, value)
	return r
}

// SetFormDatas 设置 form-data 请求体参数
func (r *Request) SetFormDatas(params map[string]string) *Request {
	for k, v := range params {
		r.SetFormData(k, v)
	}
	return r
}

// SetBody 设置请求Body
func (r *Request) SetBody(body []byte) *Request {
	r.clock.Lock()
	defer r.clock.Unlock()
	r.Body = body
	return r
}

// SetBodyString 设置请求Body
func (r *Request) SetBodyString(body string) *Request {
	r.clock.Lock()
	defer r.clock.Unlock()
	r.Body = []byte(body)
	return r
}

// SetBasicAuth 配置 Basic 认证
func (r *Request) SetBasicAuth(username, password string) *Request {
	r.clock.Lock()
	defer r.clock.Unlock()
	if r.BasicAuth == nil {
		r.BasicAuth = &BasicAuth{}
	}
	r.BasicAuth.Username = username
	r.BasicAuth.Password = password
	return r
}

// AllowRedirect 允许重定向
func (r *Request) AllowRedirect() *Request {
	r.clock.Lock()
	defer r.clock.Unlock()
	r.allowRedirect = true
	return r
}

// AllowSaveResponseHistory 允许保存重定向历史响应记录
func (r *Request) AllowSaveResponseHistory() *Request {
	r.clock.Lock()
	defer r.clock.Unlock()
	r.allowSaveResponseHistory = true
	return r
}

// SetMaxRedirectsCount 设置最大重定向请求次数支持
func (r *Request) SetMaxRedirectsCount(max int) *Request {
	r.clock.Lock()
	defer r.clock.Unlock()
	r.maxRedirectsCount = max
	return r
}

// String 用于输出完整的流量包格式
func (r *Request) String() string {
	return r.OriginalRequest.String()
}
