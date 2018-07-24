package main

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	//"database/sql"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/tealeg/xlsx"

	//_ "github.com/go-sql-driver/mysql"

	"github.com/iakud/crawler"
)

func createClient(proxy string) *crawler.Client {
	client := crawler.NewClient()
	client.Header.Add("Connection", "keep-alive")
	client.Header.Add("Upgrade-Insecure-Requests", "1")
	client.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36")
	client.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	// client.Header.Add("Accept-Encoding", "gzip, deflate")
	client.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")
	if len(proxy) > 0 {
		client.SetProxy(proxy)
	}
	// cache cookie
	client.Head("http://www.meituan.com/")
	return client
}

func main() {
	client := createClient("")

	cityMap, err := GetCityMap(client)
	if err != nil {
		log.Fatalln(err)
	}
	city, ok := cityMap["上海"]
	if !ok {
		log.Fatalln("city not found")
	}
	log.Println("city:", city.Name)

	poiIdMap, err := GetCityMeishiPoiIdMap(client, city.Acronym)
	if err != nil {
		log.Fatalln(err)
	}

	filename := fmt.Sprintf("%s_meishi_info.txt", city.Acronym)
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	xlsxFile := xlsx.NewFile()
	sheet, err := xlsxFile.AddSheet("美食")
	if err != nil {
		log.Fatalln(err)
	}
	row := sheet.AddRow()
	title := []string{"商户Id", "商户名称", "评分", "地址", "电话", "食品安全档案", "人均", "品牌Id", "品牌名称", "显示状态", "美食"}
	row.WriteSlice(&title, len(title))
	rd := bufio.NewReader(file)
	for {
		line, _, err := rd.ReadLine()
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

		row := sheet.AddRow()
		row.AddCell().SetInt64(meishiInfo.PoiId)
		row.AddCell().SetString(meishiInfo.Name)
		row.AddCell().SetFloat(float64(meishiInfo.AvgScore))
		row.AddCell().SetString(meishiInfo.Address)
		row.AddCell().SetString(meishiInfo.Phone)
		row.AddCell().SetBool(meishiInfo.HasFoodSafeInfo)
		row.AddCell().SetInt(meishiInfo.AvgPrice)
		row.AddCell().SetInt(meishiInfo.BrandId)
		row.AddCell().SetString(meishiInfo.BrandName)
		row.AddCell().SetInt(meishiInfo.ShowStatus)
		row.AddCell().SetBool(meishiInfo.IsMeishi)
	}
	xlsxFile.Save("meishi.xlsx")
	log.Println("poi num:", len(poiIdMap))

	wr := bufio.NewWriter(file)
	defer wr.Flush()

	proxy_client := createClient("//127.0.0.1:1080")
	count := 0
	totalCount := 0
	for len(poiIdMap) > 0 {
		for poiId, _ := range poiIdMap {
			meishiInfo, err := GetMeishiInfo(client, poiId)
			if err != nil {
				log.Println(err, totalCount)
				// time.Sleep(time.Second * 5)
				continue
			}
			data, err := json.Marshal(meishiInfo)
			if err != nil {
				log.Println(err, totalCount)
				// time.Sleep(time.Second * 5)
				continue
			}
			wr.Write(data)
			wr.WriteByte('\n')
			wr.Flush()
			delete(poiIdMap, poiId)
			count++
			totalCount++
			log.Println(meishiInfo.PoiId, meishiInfo.Name, meishiInfo.Phone)
			time.Sleep(time.Millisecond*100 + time.Millisecond*time.Duration(rand.Int31n(100)))
			client, proxy_client = proxy_client, client // swap
			if count > 200 {
				log.Println("reset all client")
				time.Sleep(time.Second * 15)
				client = createClient("")
				proxy_client = createClient("//127.0.0.1:1080")
				count = 0
			}
		}
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
