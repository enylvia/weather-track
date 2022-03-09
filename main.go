package main

import (
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"path/filepath"
	"weathertrack/handler"
)

func main() {
	router := gin.Default()
	webHandler := handler.NewWeatherHandler()
	router.LoadHTMLGlob("web/templates/**/*")
	router.HTMLRender = loadTemplates("./web/templates")

	router.Static("/css", "./web/assets/css")
	router.Static("/fonts", "./web/assets/fonts")
	router.Static("/img", "./web/assets/img")
	router.Static("/js", "./web/assets/js")

	router.GET("/", webHandler.Index)
	router.POST("/weather", webHandler.GetWeather)
	router.GET("/api/weather", webHandler.GetWeatherAPI)

	router.Run(":8080")

}

func loadTemplates(templatesDir string) multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	layouts, err := filepath.Glob(templatesDir + "/layouts/*")
	if err != nil {
		panic(err.Error())
	}

	includes, err := filepath.Glob(templatesDir + "/**/*")
	if err != nil {
		panic(err.Error())
	}

	// Generate our templates map from our layouts/ and includes/ directories
	for _, include := range includes {
		layoutCopy := make([]string, len(layouts))
		copy(layoutCopy, layouts)
		files := append(layoutCopy, include)
		r.AddFromFiles(filepath.Base(include), files...)
	}
	return r
}
