package crawler

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"golang.org/x/net/html"
)

type Client struct {
	httpClient *http.Client
	Header     http.Header
}

func NewClient() *Client {
	client := &Client{
		httpClient: newHttpClient(),
		Header:     make(http.Header),
	}
	return client
}

// 设置代理
func (this *Client) SetProxy(proxyURL string) error {
	u, err := url.Parse(proxyURL)
	if err != nil {
		return err
	}
	transport := &http.Transport{
		Proxy: http.ProxyURL(u),
	}
	this.httpClient.Transport = transport
	return nil
}

func (this *Client) Head(url string) error {
	_, err := this.httpClient.Head(url)
	if err != nil {
		return err
	}
	return nil
}

func (this *Client) Get(url string) (*Document, error) {
	var err error
	var req *http.Request
	req, err = this.newRequest(http.MethodGet, url)
	if err != nil {
		return nil, err
	}
	var resp *http.Response
	resp, err = this.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	node, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}
	return &Document{root: node}, nil
}

func (this *Client) newRequest(method string, url string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	this.copyHeaders(req)
	return req, nil
}

func (this *Client) copyHeaders(req *http.Request) {
	for key, values := range this.Header {
		req.Header[key] = values
	}
}

func newHttpClient() *http.Client {
	httpClient := new(http.Client)
	if jar, err := cookiejar.New(nil); err == nil {
		httpClient.Jar = jar
	}
	return httpClient
}
