package main

type City struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Acronym string `json:"acronym"`
}

type CityList []*City

func (this CityList) GetCityByName(name string) *City {
	for _, city := range this {
		if city.Name == name {
			return city
		}
	}
	return nil
}
