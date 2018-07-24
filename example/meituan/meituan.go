package main

import (
	"encoding/json"
	"errors"

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

func (this *Meituan) GetCityList() (CityList, error) {
	url := "http://www.meituan.com/changecity/"
	document, err := this.client.Get(url)
	if err != nil {
		return nil, err
	}
	datas, ok := document.Find("window.AppData = (.*);")
	if !ok {
		return nil, errors.New("city list not found")
	}
	appData := &struct {
		OpenCityList [][]interface{} `json:"openCityList"`
	}{}
	if err := json.Unmarshal([]byte(datas[0]), appData); err != nil {
		return nil, err
	}
	var cityList CityList
	for _, openCity := range appData.OpenCityList {
		if len(openCity) < 2 {
			return nil, errors.New("city list not found")
		}
		data, err := json.Marshal(openCity[1])
		if err != nil {
			return nil, err
		}
		var citys []*City
		if err := json.Unmarshal(data, &citys); err != nil {
			return nil, err
		}
		cityList = append(cityList, citys...)
	}
	return cityList, nil
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
