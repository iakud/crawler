package main

import (
	"encoding/json"
	"log"

	"github.com/iakud/crawler"
)

type City struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Acronym string `json:"acronym"`
}

type CityMap map[string]*City

func GetCityMap(client *crawler.Client) CityMap {
	url := "http://www.meituan.com/changecity/"
	document, err := client.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	datas, ok := document.Find("window.AppData = (.*);")
	if !ok {
		log.Fatalln("filters not found")
	}
	cityAppData := &struct {
		OpenCityList [][]interface{} `json:"openCityList"`
	}{}
	if err := json.Unmarshal([]byte(datas[0]), cityAppData); err != nil {
		log.Fatalln(err)
	}
	cityMap := make(CityMap)
	for _, openCity := range cityAppData.OpenCityList {
		if len(openCity) < 2 {
			continue
		}
		data, err := json.Marshal(openCity[1])
		if err != nil {
			log.Fatalln(err)
		}
		var citys []*City
		if err := json.Unmarshal(data, &citys); err != nil {
			log.Fatalln(err)
		}
		for _, city := range citys {
			cityMap[city.Name] = city
		}
	}
	return cityMap
}
