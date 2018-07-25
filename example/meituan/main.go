package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	meituan := NewMeituan()
	meituan.HeadMainPage()
	cityList, err := meituan.GetCityList()
	if err != nil {
		log.Fatalln(err)
	}
	city := cityList.GetCity("上海")
	if city == nil {
		log.Fatalln("city not found")
	}
	fmt.Println(city.Id, city.Name, city.Acronym)
	if err := meituan.WalkMeishi(city.Acronym, func(poiInfos []*MeishiPoiInfoData) {
		for _, poiInfo := range poiInfos {
			detailInfo, err := meituan.GetMeishiDetailInfo(poiInfo.PoiId)
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Println(detailInfo.PoiId, detailInfo.Name, detailInfo.Phone)
			time.Sleep(time.Second)
		}
	}); err != nil {
		log.Fatalln(err)
	}
}
