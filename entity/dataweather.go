package entity

import (
	"math"
	"strings"
)

type DataWeather struct {
	Name string `json:"name"`
	List []struct {
		Main struct {
			Temp float64 `json:"temp"`
			FeelsLike float64 `json:"feels_like"`
			Humidity int `json:"humidity"`
		} `json:"main"`
		Wind struct {
			Speed float64 `json:"speed"`
		} `json:"wind"`
		Weather []struct {
			Main        string `json:"main"`
			Description string `json:"description"`
			Icon	string `json:"icon"`
		} `json:"weather"`
		Dt string `json:"dt_txt"`
	} `json:"list"`
	City struct {
		Name       string `json:"name"`
		Country    string `json:"country"`
		Population int    `json:"population"`
	} `json:"city"`
}

type FormatResponseWeather struct {
	City struct {
		Name       string `json:"name"`
		Country    string `json:"country"`
		Population int    `json:"population"`
		Weather    []struct {
			Temp float64 `json:"temp"`
			Feels float64 `json:"feels_like"`
			Humidity int `json:"humidity"`
			Wind float64 `json:"wind"`
			Main string `json:"main"`
			Description string `json:"description"`
			Date        string `json:"date"`
			Time        string `json:"time"`
			Icon string `json:"icon"`
		} `json:"weather"`
	} `json:"city"`
}

// Set Format Response to struct
func WeatherNewFormat(dataweath DataWeather) FormatResponseWeather {
	var format FormatResponseWeather
	format.City.Name = dataweath.City.Name
	format.City.Country = dataweath.City.Country
	format.City.Population = dataweath.City.Population
	for _, val := range dataweath.List[:] {
		var weather struct {
			Temp        float64 `json:"temp"`
			Feels		float64 `json:"feels_like"`
			Humidity	int `json:"humidity"`
			Wind        float64 `json:"wind"`
			Main string `json:"main"`
			Description string  `json:"description"`
			Date        string  `json:"date"`
			Time        string  `json:"time"`
			Icon string `json:"icon"`
		}
		weather.Temp = math.Round(val.Main.Temp - 273.15)
		weather.Feels = math.Round(val.Main.FeelsLike - 273.15)
		weather.Humidity = val.Main.Humidity
		weather.Wind = math.Round(val.Wind.Speed * 3.6)
		weather.Main = val.Weather[0].Main
		weather.Description = strings.Title(val.Weather[0].Description)
		weather.Date = strings.SplitN(val.Dt, " ", 2)[0]
		weather.Time = strings.SplitN(val.Dt, " ", 2)[1]
		weather.Icon = val.Weather[0].Icon
		format.City.Weather = append(format.City.Weather, weather)
	}
	return format
}
