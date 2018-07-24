package main

import (
	"fmt"
	"log"
)

func main() {
	meituan := NewMeituan()
	cityList, err := meituan.GetCityList()
	if err != nil {
		log.Fatalln(err)
	}
	names := make(map[string]struct{})
	for _, city := range cityList {
		if _, ok := names[city.Name]; ok {
			log.Fatalln("repeat city")
		}
		fmt.Println(city.Id, city.Name, city.Acronym)
	}
}
