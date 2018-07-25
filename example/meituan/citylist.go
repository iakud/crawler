package main

import (
	"encoding/json"
	"errors"
)

type CityData struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Acronym string `json:"acronym"`
}

type CityList struct {
	citys []*CityData
}

func (this CityList) GetCity(name string) *CityData {
	for _, city := range this.citys {
		if city.Name == name {
			return city
		}
	}
	return nil
}

func (this CityList) GetCitys() []*CityData {
	var citys []*CityData
	citys = append(citys, this.citys...)
	return citys
}

func (this *Meituan) GetCityList() (*CityList, error) {
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
	var citys []*CityData
	for _, openCity := range appData.OpenCityList {
		if len(openCity) < 2 {
			return nil, errors.New("city list not found")
		}
		data, err := json.Marshal(openCity[1])
		if err != nil {
			return nil, err
		}
		var letterCitys []*CityData
		if err := json.Unmarshal(data, &letterCitys); err != nil {
			return nil, err
		}
		citys = append(citys, letterCitys...)
	}
	return &CityList{citys}, nil
}
