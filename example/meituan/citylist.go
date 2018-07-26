package main

import (
	"encoding/gob"
	"encoding/json"
	"errors"
	"os"
)

type CityData struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Acronym string `json:"acronym"`
}

type CityList struct {
	citys []*CityData
}

func (this *CityList) GetCity(name string) *CityData {
	for _, city := range this.citys {
		if city.Name == name {
			return city
		}
	}
	return nil
}

func (this *CityList) GetCitys() []*CityData {
	var citys []*CityData
	citys = append(citys, this.citys...)
	return citys
}

func (this *CityList) Save(filename string) error {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	w := gob.NewEncoder(file)
	return w.Encode(&this.citys)
}

func LoadCityList(filename string) (*CityList, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	r := gob.NewDecoder(file)
	var citys []*CityData
	if err := r.Decode(&citys); err != nil {
		return nil, err
	}
	return &CityList{citys}, nil
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
