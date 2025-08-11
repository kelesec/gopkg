package httpx

// 请求方法
const (
	MethodGet     = "GET"
	MethodHead    = "HEAD"
	MethodPost    = "POST"
	MethodPut     = "PUT"
	MethodPatch   = "PATCH"
	MethodDelete  = "DELETE"
	MethodConnect = "CONNECT"
	MethodOptions = "OPTIONS"
	MethodTrace   = "TRACE"
)

// Content-Type
const (
	MIMEApplicationJSON            = "application/json"
	MIMEApplicationJSONCharsetUTF8 = "application/json; charset=utf-8"
	MIMEApplicationXML             = "application/xml"
	MIMETextXML                    = "text/xml"
	MIMEApplicationForm            = "application/x-www-form-urlencoded"
	MIMEMultipartForm              = "multipart/form-data"
	MIMETextPlain                  = "text/plain"
	MIMETextHTML                   = "text/html"
	MIMEOctetStream                = "application/octet-stream"
)

// 一些常用的浏览器UA，需要手动添加使用
const (
	// DefaultUserAgent 默认UA
	DefaultUserAgent = "gopkg-httpx/1.0 (fasthttp-based; +https://github.com/kelesec/gopkg)"
	ChromeUserAgent  = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.5735.199 Safari/537.36"
	FirefoxUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:114.0) Gecko/20100101 Firefox/114.0"
	SafariUserAgent  = "Mozilla/5.0 (Macintosh; Intel Mac OS X 13_4_1) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.5.1 Safari/605.1.15"
	AndroidUserAgent = "Mozilla/5.0 (Linux; Android 13; SM-G998B) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/21.0 Chrome/114.0.5735.199 Mobile Safari/537.36"
)
