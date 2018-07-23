package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/iakud/crawler"
)

func main() {
	client := crawler.NewClient()
	client.Header.Add("Connection", "keep-alive")
	client.Header.Add("Upgrade-Insecure-Requests", "1")
	client.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36")
	client.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	// client.Header.Add("Accept-Encoding", "gzip, deflate")
	client.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")

	// cache cookie
	client.Head("http://www.meituan.com/")
	cityMap := GetCityMap(client)
	city, ok := cityMap["上海"]
	if !ok {
		log.Fatalln("city not found")
	}
	poiIdMap := GetCityMeishiPoiIdMap(client, city.Acronym)

	log.Println("poi num:", len(poiIdMap))
	for poiId, _ := range poiIdMap {
		meishiInfo := GetMeishiInfo(client, poiId)
		DBSave(meishiInfo)
		log.Println(meishiInfo.PoiId, meishiInfo.Name, meishiInfo.Phone)
		time.Sleep(time.Millisecond * 100)
	}
}

var defaultDb *sql.DB

func DBSave(meishiInfo *MeishiInfo) {
	if defaultDb == nil {
		var err error
		defaultDb, err = sql.Open("mysql", "kaikai:123456@tcp(192.168.2.251:3306)/kaikai?charset=utf8")
		if err != nil {
			log.Fatalln(err)
			return
		}
	}

	if meishiInfo == nil {
		return
	}
	sqlcmd := fmt.Sprintf("insert into meishi values(?,?,?,?,?,?,?,?,?,?,?)")
	poiId := meishiInfo.PoiId
	name := meishiInfo.Name
	avgScore := meishiInfo.AvgScore
	address := meishiInfo.Address
	phone := meishiInfo.Phone
	hasFoodSafeInfo := meishiInfo.HasFoodSafeInfo
	avgPrice := meishiInfo.AvgPrice
	brandId := meishiInfo.BrandId
	brandName := meishiInfo.BrandName
	showStatus := meishiInfo.ShowStatus
	isMeishi := meishiInfo.IsMeishi
	_, err := defaultDb.Exec(sqlcmd, poiId, name, avgScore, address, phone, hasFoodSafeInfo, avgPrice, brandId, brandName, showStatus, isMeishi)
	if err != nil {
		log.Println(err)
		return
	}
}

func CloseDb() {
	if defaultDb != nil {
		defaultDb.Close()
		defaultDb = nil
	}
}