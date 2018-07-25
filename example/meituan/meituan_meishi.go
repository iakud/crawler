package main

import (
	"encoding/json"
	"fmt"
)

type MeishiCateData struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Url  string `json:"url"`
}

type MeishiSubAreaData struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Url  string `json:"url"`
}

type MeishiAreaData struct {
	Id       int                  `json:"id"`
	Name     string               `json:"name"`
	Url      string               `json:"url"`
	SubAreas []*MeishiSubAreaData `json:"subAreas"`
}

type MeishiFilters struct {
	Cates []*MeishiCateData `json:"cates"`
	Areas []*MeishiAreaData `json:"areas"`
}

type MeishiPoiInfoData struct {
	PoiId         int64   `json:"poiId"`
	Title         string  `json:"title"`
	AvgScore      float64 `json:"avgScore"`
	AllCommentNum int     `json:"allCommentNum"`
	Address       string  `json:"address"`
	AvgPrice      int     `json:"avgPrice"`
}

type MeishiPoiLists struct {
	TotalCounts int                  `json:"totalCounts"`
	PoiInfos    []*MeishiPoiInfoData `json:"poiInfos"`
}

type MeishiData struct {
	Filters  *MeishiFilters  `json:"filters"`
	CateId   int             `json:"cateId"`
	AreaId   int             `json:"areaId"`
	Pn       int             `json:"pn"`
	PoiLists *MeishiPoiLists `json:"poiLists"`
}

func (this *Meituan) WalkMeishi(acronym string, walkFunc func([]*MeishiPoiInfoData)) error {
	if walkFunc == nil {
		return nil
	}
	url := fmt.Sprintf("http://%v.meituan.com/meishi/", acronym)
	data, err := this.getMeishiData(url)
	if err != nil {
		return err
	}
	if data.PoiLists.TotalCounts < 1000 || len(data.Filters.Cates) == 0 {
		if len(data.PoiLists.PoiInfos) == 0 {
			return nil
		}
		walkFunc(data.PoiLists.PoiInfos)
		if len(data.PoiLists.PoiInfos) < 32 {
			return nil
		}
		return this.walkMeishiPn(url, walkFunc)
	}
	for _, cate := range data.Filters.Cates {
		if err := this.walkMeishiCate(cate, walkFunc); err != nil {
			return err
		}
	}
	return nil
}

func (this *Meituan) walkMeishiCate(cate *MeishiCateData, walkFunc func([]*MeishiPoiInfoData)) error {
	data, err := this.getMeishiData(cate.Url)
	if err != nil {
		return err
	}
	if data.PoiLists.TotalCounts < 1000 || len(data.Filters.Areas) == 0 {
		if len(data.PoiLists.PoiInfos) == 0 {
			return nil
		}
		walkFunc(data.PoiLists.PoiInfos)
		if len(data.PoiLists.PoiInfos) < 32 {
			return nil
		}
		return this.walkMeishiPn(cate.Url, walkFunc)
	}
	for _, area := range data.Filters.Areas {
		if err := this.walkMeishiArea(area, walkFunc); err != nil {
			return err
		}
	}
	return nil
}

func (this *Meituan) walkMeishiArea(area *MeishiAreaData, walkFunc func([]*MeishiPoiInfoData)) error {
	data, err := this.getMeishiData(area.Url)
	if err != nil {
		return err
	}
	var subAreas []*MeishiSubAreaData
	for _, subArea := range area.SubAreas {
		if subArea.Id != area.Id {
			subAreas = append(subAreas, subArea)
		}
	}
	if data.PoiLists.TotalCounts < 1000 || len(subAreas) == 0 {
		if len(data.PoiLists.PoiInfos) == 0 {
			return nil
		}
		walkFunc(data.PoiLists.PoiInfos)
		if len(data.PoiLists.PoiInfos) < 32 {
			return nil
		}
		return this.walkMeishiPn(area.Url, walkFunc)
	}
	for _, subArea := range subAreas {
		if err := this.walkMeishiSubArea(subArea, walkFunc); err != nil {
			return err
		}
	}
	return nil
}

func (this *Meituan) walkMeishiSubArea(subArea *MeishiSubAreaData, walkFunc func([]*MeishiPoiInfoData)) error {
	data, err := this.getMeishiData(subArea.Url)
	if err != nil {
		return err
	}
	if len(data.PoiLists.PoiInfos) == 0 {
		return nil
	}
	walkFunc(data.PoiLists.PoiInfos)
	if len(data.PoiLists.PoiInfos) < 32 {
		return nil
	}
	return this.walkMeishiPn(subArea.Url, walkFunc)
}

func (this *Meituan) walkMeishiPn(url string, walkFunc func([]*MeishiPoiInfoData)) error {
	for pn := 2; pn <= 32; pn++ {
		pnUrl := fmt.Sprintf("%vpn%v/", url, pn)
		data, err := this.getMeishiData(pnUrl)
		if err != nil {
			return err
		}
		if len(data.PoiLists.PoiInfos) == 0 {
			return nil
		}
		walkFunc(data.PoiLists.PoiInfos)
		if len(data.PoiLists.PoiInfos) < 32 {
			return nil
		}
	}
	return nil
}

func (this *Meituan) getMeishiData(url string) (*MeishiData, error) {
	document, err := this.client.Get(url)
	if err != nil {
		return nil, err
	}
	datas, ok := document.Find("window._appState = (.*);")
	if !ok {
		return nil, fmt.Errorf("meishi list not found, url=%v", url)
	}
	appData := &MeishiData{}
	if err := json.Unmarshal([]byte(datas[0]), appData); err != nil {
		return nil, err
	}
	return appData, nil
}

type MeishiDetailInfoData struct {
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

func (this *Meituan) GetMeishiDetailInfo(poiId int64) (*MeishiDetailInfoData, error) {
	url := fmt.Sprintf("http://www.meituan.com/meishi/%v/", poiId)
	document, err := this.client.Get(url)
	if err != nil {
		return nil, err
	}
	datas, ok := document.Find("window._appState = (.*);")
	if !ok {
		return nil, fmt.Errorf("meishi not found, url=%v", url)
	}
	appData := &struct {
		DetailInfo *MeishiDetailInfoData `json:"detailInfo"`
	}{}
	if err := json.Unmarshal([]byte(datas[0]), appData); err != nil {
		return nil, err
	}
	return appData.DetailInfo, nil
}
