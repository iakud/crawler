package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/iakud/crawler"
)

type CityMeishiPoiIdMap map[int64]struct{}

func GetCityMeishiPoiIdMap(client *crawler.Client, acronym string) CityMeishiPoiIdMap {
	filename := fmt.Sprintf("%s_meishi_poiid.txt", acronym)
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatalln(err)
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
		return poiIdMap
	}

	url := fmt.Sprintf("http://%v.meituan.com/meishi/", acronym)
	cityMeishi := getCityMeishi(client, url)

	if cityMeishi.PoiLists.TotalCounts < 1000 {
		pagePoiIdMap := getCityMeishiPagePoiIdMap(client, url)
		for poiId, _ := range pagePoiIdMap {
			if _, ok := poiIdMap[poiId]; !ok {
				poiIdMap[poiId] = struct{}{}
			}
		}
	} else {
		for _, cate := range cityMeishi.Filters.Cates {
			catePoiIdMap := getCityMeishiCatePoiIdMap(client, cate.Url)
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
	return poiIdMap
}

func getCityMeishiCatePoiIdMap(client *crawler.Client, url string) CityMeishiPoiIdMap {
	cityMeishi := getCityMeishi(client, url)
	poiIdMap := make(CityMeishiPoiIdMap)
	if cityMeishi.PoiLists.TotalCounts < 1000 {
		pagePoiIdMap := getCityMeishiPagePoiIdMap(client, url)
		for poiId, _ := range pagePoiIdMap {
			if _, ok := poiIdMap[poiId]; !ok {
				poiIdMap[poiId] = struct{}{}
			}
		}
	} else {
		for _, area := range cityMeishi.Filters.Areas {
			areaCityMeishi := getCityMeishi(client, area.Url)
			if areaCityMeishi.PoiLists.TotalCounts < 1000 {
				pagePoiIdMap := getCityMeishiPagePoiIdMap(client, area.Url)
				for poiId, _ := range pagePoiIdMap {
					if _, ok := poiIdMap[poiId]; !ok {
						poiIdMap[poiId] = struct{}{}
					}
				}
				log.Println(area.Name, len(pagePoiIdMap))
			} else {
				for _, subArea := range area.SubAreas {
					pagePoiIdMap := getCityMeishiPagePoiIdMap(client, subArea.Url)
					for poiId, _ := range pagePoiIdMap {
						if _, ok := poiIdMap[poiId]; !ok {
							poiIdMap[poiId] = struct{}{}
						}
					}
				}
			}
		}
	}
	return poiIdMap
}

func getCityMeishiPagePoiIdMap(client *crawler.Client, url string) CityMeishiPoiIdMap {
	poiIdMap := make(CityMeishiPoiIdMap)
	for i := 1; i <= 32; i++ {
		pageUrl := fmt.Sprintf("%spn%v/", url, i)
		cityMeishi := getCityMeishi(client, pageUrl)
		if len(cityMeishi.PoiLists.PoiInfos) == 0 {
			break
		}
		for _, poiInfo := range cityMeishi.PoiLists.PoiInfos {
			poiIdMap[poiInfo.PoiId] = struct{}{}
		}
	}
	return poiIdMap
}

func getCityMeishi(client *crawler.Client, url string) *CityMeishi {
	document, err := client.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	datas, ok := document.Find("window._appState = (.*);")
	if !ok {
		log.Fatalln("city meishi not found", url)
	}
	cityMeishi := &CityMeishi{}
	if err := json.Unmarshal([]byte(datas[0]), cityMeishi); err != nil {
		log.Fatalln(err)
	}
	return cityMeishi
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

func GetMeishiInfo(client *crawler.Client, poiId int64) *MeishiInfo {
	url := fmt.Sprintf("http://www.meituan.com/meishi/%v/", poiId)
	document, err := client.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	datas, ok := document.Find("window._appState = (.*);")
	if !ok {
		log.Fatalln("city meishi not found", url)
	}
	meishiAppState := &struct {
		DetailInfo *MeishiInfo `json:"detailInfo"`
	}{}
	if err := json.Unmarshal([]byte(datas[0]), meishiAppState); err != nil {
		log.Fatalln(err)
	}
	return meishiAppState.DetailInfo
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
