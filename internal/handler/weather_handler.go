package handler

import (
	"api-clima-go/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WeatherHandler struct {
	Service *service.WeatherService
}

func NewWeatherHandler(service *service.WeatherService) *WeatherHandler{
	return &WeatherHandler{Service: service}
}

func (h *WeatherHandler) GetWeather(c *gin.Context) {
	city := c.Param("city")

	if city == "" {
		c.JSON(http.StatusBadRequest, gin.H {
			"error": "cidade não informada!",
		})
		return
	}

	result, err := h.Service.GetFastest(city)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Cidade não informada!",
		})
		return
	}
	c.JSON(http.StatusOK, result)
}