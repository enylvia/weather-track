package handler

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"weathertrack/entity"
)

type weatherHandler struct {
}
type apiConfig struct {
	ApiKeyWeatherMap string `json:"ApiKeyWeatherMap"`
}

func NewWeatherHandler() *weatherHandler {
	return &weatherHandler{}
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

func query(city string) (entity.FormatResponseWeather, error) {
	apiConfigs, err := loadApiConfig("apiConfig.json")
	if err != nil {
		return entity.FormatResponseWeather{}, err
	}
	url := "https://api.openweathermap.org/data/2.5/forecast?APPID=" + apiConfigs.ApiKeyWeatherMap + "&q=" + city + "&lang=eng&cnt=6"

	resp, err := http.Get(url)
	if err != nil {
		return entity.FormatResponseWeather{}, errors.New("City Not Found")
	}
	defer resp.Body.Close()

	var data entity.DataWeather
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return entity.FormatResponseWeather{}, err
	}
	return entity.WeatherNewFormat(data), nil
}

func (h *weatherHandler) Index(c *gin.Context) {
	c.HTML(http.StatusOK, "getweather.html", nil)
}

func (h *weatherHandler) GetWeather(c *gin.Context) {
	city := c.PostForm("city")
	data, err := query(city)
	if err != nil {
		http.Error(c.Writer, "City Not Found", http.StatusNotFound)
		return
	}
	c.HTML(http.StatusOK, "weather.html", gin.H{"data": data})
}

func (h *weatherHandler) GetWeatherAPI(c *gin.Context) {
	city := c.Query("city")
	data, err := query(city)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}
