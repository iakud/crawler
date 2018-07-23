package crawler

import (
	"log"
	"net/http"
	"net/url"

	"golang.org/x/net/html"
)

type Client struct {
	httpClient *http.Client

	Header  http.Header
	cookies map[string]*http.Cookie
}

func NewClient() *Client {
	client := &Client{
		httpClient: &http.Client{},

		Header:  make(http.Header),
		cookies: make(map[string]*http.Cookie),
	}
	return client
}

// 设置代理
func (this *Client) SetProxy(rawurl string) {
	u, err := url.Parse(rawurl)
	if err != nil {
		log.Println(err)
		return
	}
	transport := &http.Transport{
		Proxy: http.ProxyURL(u),
	}
	this.httpClient.Transport = transport
}

func (this *Client) Head(url string) error {
	resp, err := this.httpClient.Head(url)
	if err != nil {
		return err
	}
	// save cookie
	for _, cookie := range resp.Cookies() {
		this.cookies[cookie.Name] = cookie
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
	// set headers
	for key, values := range this.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
	// set cookies
	for _, cookie := range this.cookies {
		req.AddCookie(cookie)
	}
	return req, nil
}
