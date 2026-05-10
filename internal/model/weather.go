package model

type WeatherResponse struct {
	City string `json:"city"`
	Temperature float64 `json:"temperature"`
	Humidity int `json:"humidity"`
	Condition string `json:"condition"`
	Source string `json:"source"`
}