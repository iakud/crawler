package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/iakud/crawler"
)

type CityMeishiPoiIdMap map[int64]struct{}

func GetCityMeishiPoiIdMap(client *crawler.Client, acronym string) (CityMeishiPoiIdMap, error) {
	filename := fmt.Sprintf("%s_meishi_poiid.txt", acronym)
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	poiIdMap := make(CityMeishiPoiIdMap)
	rd := bufio.NewReader(file)
	for {
		line, _, err := rd.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalln(err)
		}
		poiId, err := strconv.ParseInt(string(line), 10, 64)
		if err != nil {
			log.Fatalln(err)
		}
		poiIdMap[poiId] = struct{}{}
	}

	if len(poiIdMap) > 0 {
		return poiIdMap, nil
	}

	url := fmt.Sprintf("http://%v.meituan.com/meishi/", acronym)
	cityMeishi, err := getCityMeishi(client, url)
	if err != nil {
		return nil, err
	}

	if cityMeishi.PoiLists.TotalCounts < 1000 {
		pagePoiIdMap, err := getCityMeishiPagePoiIdMap(client, url)
		if err != nil {
			return nil, err
		}
		for poiId, _ := range pagePoiIdMap {
			if _, ok := poiIdMap[poiId]; !ok {
				poiIdMap[poiId] = struct{}{}
			}
		}
	} else {
		for _, cate := range cityMeishi.Filters.Cates {
			catePoiIdMap, err := getCityMeishiCatePoiIdMap(client, cate.Url)
			if err != nil {
				return nil, err
			}
			for poiId, _ := range catePoiIdMap {
				if _, ok := poiIdMap[poiId]; !ok {
					poiIdMap[poiId] = struct{}{}
				}
			}
			log.Println(cate.Name, len(catePoiIdMap))
		}
	}
	wr := bufio.NewWriter(file)
	defer wr.Flush()
	for poiId, _ := range poiIdMap {
		wr.WriteString(strconv.FormatInt(poiId, 10))
		wr.WriteByte('\n')
	}
	return poiIdMap, nil
}

func getCityMeishiCatePoiIdMap(client *crawler.Client, url string) (CityMeishiPoiIdMap, error) {
	cityMeishi, err := getCityMeishi(client, url)
	if err != nil {
		return nil, err
	}
	poiIdMap := make(CityMeishiPoiIdMap)
	if cityMeishi.PoiLists.TotalCounts < 1000 {
		pagePoiIdMap, err := getCityMeishiPagePoiIdMap(client, url)
		if err != nil {
			return nil, err
		}
		for poiId, _ := range pagePoiIdMap {
			if _, ok := poiIdMap[poiId]; !ok {
				poiIdMap[poiId] = struct{}{}
			}
		}
	} else {
		for _, area := range cityMeishi.Filters.Areas {
			areaCityMeishi, err := getCityMeishi(client, area.Url)
			if err != nil {
				return nil, err
			}
			if areaCityMeishi.PoiLists.TotalCounts < 1000 {
				pagePoiIdMap, err := getCityMeishiPagePoiIdMap(client, area.Url)
				if err != nil {
					return nil, err
				}
				for poiId, _ := range pagePoiIdMap {
					if _, ok := poiIdMap[poiId]; !ok {
						poiIdMap[poiId] = struct{}{}
					}
				}
				log.Println(area.Name, len(pagePoiIdMap))
			} else {
				for _, subArea := range area.SubAreas {
					pagePoiIdMap, err := getCityMeishiPagePoiIdMap(client, subArea.Url)
					if err != nil {
						return nil, err
					}
					for poiId, _ := range pagePoiIdMap {
						if _, ok := poiIdMap[poiId]; !ok {
							poiIdMap[poiId] = struct{}{}
						}
					}
				}
			}
		}
	}
	return poiIdMap, nil
}

func getCityMeishiPagePoiIdMap(client *crawler.Client, url string) (CityMeishiPoiIdMap, error) {
	poiIdMap := make(CityMeishiPoiIdMap)
	for i := 1; i <= 32; i++ {
		pageUrl := fmt.Sprintf("%spn%v/", url, i)
		cityMeishi, err := getCityMeishi(client, pageUrl)
		if err != nil {
			return nil, err
		}
		if len(cityMeishi.PoiLists.PoiInfos) == 0 {
			break
		}
		for _, poiInfo := range cityMeishi.PoiLists.PoiInfos {
			poiIdMap[poiInfo.PoiId] = struct{}{}
		}
	}
	return poiIdMap, nil
}

func getCityMeishi(client *crawler.Client, url string) (*CityMeishi, error) {
	document, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	datas, ok := document.Find("window._appState = (.*);")
	if !ok {
		return nil, errors.New("city meishi not found")
	}
	cityMeishi := &CityMeishi{}
	if err := json.Unmarshal([]byte(datas[0]), cityMeishi); err != nil {
		return nil, err
	}
	return cityMeishi, nil
}

// 分类
type CityMeishi struct {
	Filters *struct {
		Cates []*struct {
			Id   int    `json:"id"`
			Name string `json:"name"`
			Url  string `json:"url"`
		} `json:"cates"`
		Areas []*struct {
			Id       int    `json:"id"`
			Name     string `json:"name"`
			Url      string `json:"url"`
			SubAreas []*struct {
				Id   int    `json:"id"`
				Name string `json:"name"`
				Url  string `json:"url"`
			} `json:"subAreas"`
		} `json:"areas"`
	} `json:"filters"`
	Pn       int `json:"pn"`
	PoiLists *struct {
		TotalCounts int `json:"totalCounts"`
		PoiInfos    []*struct {
			PoiId int64  `json:"poiId"`
			Title string `json:"title"`
		} `json:"poiInfos"`
	} `json:"poiLists"`
}

func GetMeishiInfo(client *crawler.Client, poiId int64) (*MeishiInfo, error) {
	url := fmt.Sprintf("http://www.meituan.com/meishi/%v/", poiId)
	document, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	datas, ok := document.Find("window._appState = (.*);")
	if !ok {
		return nil, fmt.Errorf("meishi not found, url=%v", url)
	}
	meishiAppState := &struct {
		DetailInfo *MeishiInfo `json:"detailInfo"`
	}{}
	if err := json.Unmarshal([]byte(datas[0]), meishiAppState); err != nil {
		return nil, err
	}
	return meishiAppState.DetailInfo, nil
}

type MeishiInfo struct {
	PoiId    int64   `json:"poiId"`
	Name     string  `json:"name"`
	AvgScore float32 `json:"avgScore"`
	Address  string  `json:"address"`
	Phone    string  `json:"phone"`

	HasFoodSafeInfo bool   `json:"hasFoodSafeInfo"`
	AvgPrice        int    `json:"avgPrice"`
	BrandId         int    `json:"brandId"`
	BrandName       string `json:"brandName"`
	ShowStatus      int    `json:"showStatus"`
	IsMeishi        bool   `json:"isMeishi"`
}
