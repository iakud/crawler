package main

import (
	"bufio"
	"fmt"
	"io"
	//"database/sql"
	"encoding/json"
	"log"
	"os"
	"time"

	//_ "github.com/go-sql-driver/mysql"

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
	log.Println("city:", city.Name)

	poiIdMap := GetCityMeishiPoiIdMap(client, city.Acronym)
	log.Println("poi num:", len(poiIdMap))

	filename := fmt.Sprintf("%s_meishi_info.txt", city.Acronym)
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()
	rd := bufio.NewReader(file)
	for {
		line, err := rd.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalln(err)
		}
		meishiInfo := &MeishiInfo{}
		if err := json.Unmarshal(line, meishiInfo); err != nil {
			log.Fatalln(err)
		}
		delete(poiIdMap, meishiInfo.PoiId)
	}
	wr := bufio.NewWriter(file)
	defer wr.Flush()
	for poiId, _ := range poiIdMap {
		meishiInfo := GetMeishiInfo(client, poiId)
		data, err := json.Marshal(meishiInfo)
		if err != nil {
			log.Fatalln(err)
		}
		wr.Write(data)
		wr.WriteByte('\n')
		log.Println(meishiInfo.PoiId, meishiInfo.Name, meishiInfo.Phone)
		time.Sleep(time.Second)
	}
}

/*
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
}*/
