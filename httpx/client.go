package httpx

import (
	"crypto/tls"
	"fmt"
	"github.com/chainreactors/proxyclient"
	"github.com/valyala/fasthttp"
	"net"
	"net/url"
	"sync"
	"time"
)

// Client 请求客户端，参数和 `fasthttp.Client` 基本一致
type Client struct {
	ReadTimeout                   time.Duration     // 响应读取超时
	WriteTimeout                  time.Duration     // 请求写入超时
	MaxIdleConnDuration           time.Duration     // 空闲连接保留时间（keep-alive 连接释放允许的最大时间）
	MaxConnWaitTimeout            time.Duration     // 连接池耗尽之后，新的连接请求超时等待时间
	ReadBufferSize                int               // 大响应缓冲（适合下载文件或大 JSON）
	WriteBufferSize               int               // 大请求缓冲（适合 POST 大文件）
	MaxResponseBodySize           int               // 单个 HTTP 响应体的最大字节数
	MaxConnsPerHost               int               // 单主机最大并发连接数（防被封禁）
	NoDefaultUserAgentHeader      bool              // 禁用默认 UA(fasthttp)，禁用后需要手动配置
	DisableHeaderNamesNormalizing bool              // 原样发送请求头，不对请求头字段名称进行规范化
	DisablePathNormalizing        bool              // 保留原始 URL 路径，不进行规范化（不处理特殊字符等），适合发送特殊字符
	Dial                          fasthttp.DialFunc // 用于建立与主机的新连接的回调
	TLSConfig                     *tls.Config       // 证书相关配置
	fastClient                    *fasthttp.Client  // 确保请求只用到一个Client就行

	// 加个锁
	clock *sync.Mutex
}

// NewClient 使用默认配置创建 Client
func NewClient() *Client {
	return &Client{
		ReadTimeout:                   time.Second * 10,
		WriteTimeout:                  time.Second * 10,
		MaxIdleConnDuration:           time.Second * 10,
		MaxConnWaitTimeout:            time.Second * 10,
		NoDefaultUserAgentHeader:      true,
		DisableHeaderNamesNormalizing: true,
		DisablePathNormalizing:        true,
		ReadBufferSize:                4 * 1024 * 1024,
		WriteBufferSize:               1 * 1024 * 1024,
		MaxResponseBodySize:           10 * 1024 * 1024,
		MaxConnsPerHost:               1024,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		Dial: (&fasthttp.TCPDialer{
			Concurrency:      4096,
			DNSCacheDuration: time.Hour,
		}).Dial,
		fastClient: &fasthttp.Client{},
		clock:      &sync.Mutex{},
	}
}

// NewClientWithFastHttp 通过 `fasthttp.Client` 进行创建
func NewClientWithFastHttp(client *fasthttp.Client) *Client {
	return &Client{
		ReadTimeout:                   client.ReadTimeout,
		WriteTimeout:                  client.WriteTimeout,
		MaxIdleConnDuration:           client.MaxIdleConnDuration,
		MaxConnWaitTimeout:            client.MaxConnWaitTimeout,
		ReadBufferSize:                client.ReadBufferSize,
		WriteBufferSize:               client.WriteBufferSize,
		MaxResponseBodySize:           client.MaxResponseBodySize,
		MaxConnsPerHost:               client.MaxConnsPerHost,
		NoDefaultUserAgentHeader:      client.NoDefaultUserAgentHeader,
		DisableHeaderNamesNormalizing: client.DisableHeaderNamesNormalizing,
		DisablePathNormalizing:        client.DisablePathNormalizing,
		Dial:                          client.Dial,
		TLSConfig:                     client.TLSConfig,
		fastClient:                    client,
		clock:                         &sync.Mutex{},
	}
}

// R 创建Request请求对象
func (cli *Client) R() *Request {
	return newRequest(cli)
}

// preCheck fastClient 前置检查，确保数据都完全同步给 fastClient
func (cli *Client) preCheck() {
	cli.clock.Lock()
	defer cli.clock.Unlock()

	if cli.fastClient == nil {
		cli.fastClient = &fasthttp.Client{}
	}

	// 将内容全部同步给 fastClient
	cli.fastClient.ReadTimeout = cli.ReadTimeout
	cli.fastClient.WriteTimeout = cli.WriteTimeout
	cli.fastClient.MaxIdleConnDuration = cli.MaxIdleConnDuration
	cli.fastClient.MaxConnWaitTimeout = cli.MaxConnWaitTimeout
	cli.fastClient.ReadBufferSize = cli.ReadBufferSize
	cli.fastClient.WriteBufferSize = cli.WriteBufferSize
	cli.fastClient.MaxResponseBodySize = cli.MaxResponseBodySize
	cli.fastClient.MaxConnsPerHost = cli.MaxConnsPerHost
	cli.fastClient.NoDefaultUserAgentHeader = cli.NoDefaultUserAgentHeader
	cli.fastClient.DisableHeaderNamesNormalizing = cli.DisableHeaderNamesNormalizing
	cli.fastClient.DisablePathNormalizing = cli.DisablePathNormalizing
	cli.fastClient.Dial = cli.Dial
	cli.fastClient.TLSConfig = cli.TLSConfig
}

// execute 执行 HTTP 请求
func (cli *Client) execute(req *fasthttp.Request, resp *fasthttp.Response) error {
	cli.preCheck()
	return cli.fastClient.Do(req, resp)
}

func (cli *Client) SetReadTimeout(t time.Duration) *Client {
	cli.clock.Lock()
	defer cli.clock.Unlock()
	cli.ReadTimeout = t
	return cli
}

func (cli *Client) SetWriteTimeout(t time.Duration) *Client {
	cli.clock.Lock()
	defer cli.clock.Unlock()
	cli.WriteTimeout = t
	return cli
}

func (cli *Client) SetMaxIdleConnDuration(t time.Duration) *Client {
	cli.clock.Lock()
	defer cli.clock.Unlock()
	cli.MaxIdleConnDuration = t
	return cli
}

func (cli *Client) SetMaxConnWaitTimeout(t time.Duration) *Client {
	cli.clock.Lock()
	defer cli.clock.Unlock()
	cli.MaxConnWaitTimeout = t
	return cli
}

func (cli *Client) SetNoDefaultUserAgentHeader(b bool) *Client {
	cli.clock.Lock()
	defer cli.clock.Unlock()
	cli.NoDefaultUserAgentHeader = b
	return cli
}

func (cli *Client) SetDisableHeaderNamesNormalizing(b bool) *Client {
	cli.clock.Lock()
	defer cli.clock.Unlock()
	cli.DisableHeaderNamesNormalizing = b
	return cli
}

func (cli *Client) SetDisablePathNormalizing(b bool) *Client {
	cli.clock.Lock()
	defer cli.clock.Unlock()
	cli.DisablePathNormalizing = b
	return cli
}

func (cli *Client) SetReadBufferSize(n int) *Client {
	cli.clock.Lock()
	defer cli.clock.Unlock()
	cli.ReadBufferSize = n
	return cli
}

func (cli *Client) SetWriteBufferSize(n int) *Client {
	cli.clock.Lock()
	defer cli.clock.Unlock()
	cli.WriteBufferSize = n
	return cli
}

func (cli *Client) SetMaxResponseBodySize(n int) *Client {
	cli.clock.Lock()
	defer cli.clock.Unlock()
	cli.MaxResponseBodySize = n
	return cli
}

func (cli *Client) SetMaxConnsPerHost(n int) *Client {
	cli.clock.Lock()
	defer cli.clock.Unlock()
	cli.MaxConnsPerHost = n
	return cli
}

func (cli *Client) SetTLSConfig(c *tls.Config) *Client {
	cli.clock.Lock()
	defer cli.clock.Unlock()
	cli.TLSConfig = c
	return cli
}

func (cli *Client) SetDial(f fasthttp.DialFunc) *Client {
	cli.clock.Lock()
	defer cli.clock.Unlock()
	cli.Dial = f
	return cli
}

// SetProxy 设置请求代理，采用 github.com/chainreactors/proxyclient 库作为代理支持
// 因此 proxyclient 库支持的协议应该都支持，如 HTTP/HTTPS/SOCKS5 等
func (cli *Client) SetProxy(proxy string) *Client {
	cli.clock.Lock()
	defer cli.clock.Unlock()

	proxyURL, err := url.Parse(proxy)
	if err != nil {
		panic(fmt.Errorf("invalid proxy url: %s", proxy))
	}

	dial, err := proxyclient.NewClient(proxyURL)
	if err != nil {
		panic(fmt.Errorf("create client failed: %s", err))
	}

	cli.Dial = func(addr string) (net.Conn, error) {
		return dial.Dial("tcp", addr)
	}

	return cli
}

// SetProxies 设置请求代理池，采用 github.com/chainreactors/proxyclient 库作为代理支持
// 因此 proxyclient 库支持的协议应该都支持，如 HTTP/HTTPS/SOCKS5 等
func (cli *Client) SetProxies(proxies []string) *Client {
	cli.clock.Lock()
	defer cli.clock.Unlock()

	var proxyUrls []*url.URL
	for _, proxy := range proxies {
		proxyURL, err := url.Parse(proxy)
		if err != nil {
			panic(fmt.Errorf("invalid proxy url: %s", proxy))
		}
		proxyUrls = append(proxyUrls, proxyURL)
	}

	dialer, err := proxyclient.NewClientChain(proxyUrls)
	if err != nil {
		panic(fmt.Errorf("create client failed: %s", err))
	}

	cli.Dial = func(addr string) (net.Conn, error) {
		return dialer.Dial("tcp", addr)
	}

	return cli
}
