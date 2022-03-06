package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strings"
)

type apiConfig struct {
	ApiKeyWeatherMap string `json:"ApiKeyWeatherMap"`
}

//Mapping response from API to struct
type DataWeather struct {
	Name string `json:"name"`
	List []struct {
		Main struct {
			Temp float64 `json:"temp"`
		} `json:"main"`
		Wind struct {
			Speed float64 `json:"speed"`
		} `json:"wind"`
		Weather []struct {
			Main        string `json:"main"`
			Description string `json:"description"`
		} `json:"weather"`
		Dt string `json:"dt_txt"`
	} `json:"list"`
	City struct {
		Name       string `json:"name"`
		Country    string `json:"country"`
		Population int    `json:"population"`
	} `json:"city"`
}

//Make New Format for response
type FormatResponseWeather struct {
	City struct {
		Name       string `json:"name"`
		Country    string `json:"country"`
		Population int    `json:"population"`
		Weather    []struct {
			Temp        float64 `json:"temp"`
			Wind        float64 `json:"wind"`
			Description string  `json:"description"`
			Date        string  `json:"date"`
			Time        string  `json:"time"`
		} `json:"weather"`
	} `json:"city"`
}
// Set Format Response to struct
func WeatherNewFormat(dataweath DataWeather) FormatResponseWeather {
	var format FormatResponseWeather
	format.City.Name = dataweath.City.Name
	format.City.Country = dataweath.City.Country
	format.City.Population = dataweath.City.Population
	for _, val := range dataweath.List[:4] {
		var weather struct {
			Temp        float64 `json:"temp"`
			Wind        float64 `json:"wind"`
			Description string  `json:"description"`
			Date        string  `json:"date"`
			Time        string  `json:"time"`
		}
		weather.Temp = math.Round(val.Main.Temp - 273.15)
		weather.Wind = val.Wind.Speed
		weather.Description = val.Weather[0].Description
		weather.Date = strings.SplitN(val.Dt, " ", 2)[0]
		weather.Time = strings.SplitN(val.Dt, " ", 2)[1]
		format.City.Weather = append(format.City.Weather, weather)
	}
	return format
}
func loadApiConfig(filename string) (apiConfig, error) {
	bytes, err := ioutil.ReadFile(filename)

	if err != nil {
		return apiConfig{}, err
	}
	var config apiConfig

	err = json.Unmarshal(bytes, &config)

	if err != nil {
		return apiConfig{}, err
	}

	return config, nil
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

func query(city string) (FormatResponseWeather, error) {
	apiConfigs, err := loadApiConfig("apiConfig.json")
	if err != nil {
		return FormatResponseWeather{}, err
	}
	url := "https://api.openweathermap.org/data/2.5/forecast?APPID=" + apiConfigs.ApiKeyWeatherMap + "&q=" + city + "&lang=id&cnt=4"

	resp, err := http.Get(url)
	if err != nil {
		return FormatResponseWeather{}, err
	}
	defer resp.Body.Close()

	var data DataWeather
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return FormatResponseWeather{}, err
	}
	return WeatherNewFormat(data), nil
}

func main() {
	http.HandleFunc("/", hello)
	http.HandleFunc("/weather/", func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path == "/weather/" {
			fmt.Printf("[ERROR] Location should not be empty")
			http.Redirect(writer, request, "/", http.StatusSeeOther)
			return
		}
		city := strings.SplitN(request.URL.Path, "/", 3)[2]
		data, err := query(city)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(writer).Encode(data)
	})
	fmt.Println("server started at localhost:9000")
	http.ListenAndServe(":8080", nil)

}
