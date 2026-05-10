package main

import (
	"api-clima-go/internal/config"
	"api-clima-go/internal/handler"
	"api-clima-go/internal/provider"
	"api-clima-go/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {

	cfg := config.Load()

	openWeather := provider.NewOpenWeatherProvider(cfg.OpenWeatherApiKey)
	weatherAPI := provider.NewWeatherAPIProvider(cfg.WeatherApiKey)

	weatherService := service.NewWeatherService(openWeather, weatherAPI)
	weatherHandler := handler.NewWeatherHandler(weatherService)

	server := gin.Default()
	server.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	server.GET("/weather/:city", weatherHandler.GetWeather)
	server.Run(":" + cfg.ServerPort)
	}