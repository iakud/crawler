package main

import (
	"github.com/iakud/crawler"
)

type Meituan struct {
	client *crawler.Client
}

func NewMeituan() *Meituan {
	meituan := &Meituan{
		client: newDefaultClient(),
	}
	return meituan
}

func (this *Meituan) HeadMainPage() {
	this.client.Head("http://www.meituan.com/")
}

func newDefaultClient() *crawler.Client {
	client := crawler.NewClient()
	client.Header.Add("Connection", "keep-alive")
	client.Header.Add("Upgrade-Insecure-Requests", "1")
	client.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36")
	client.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	// client.Header.Add("Accept-Encoding", "gzip, deflate")
	client.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")

	return client
}
