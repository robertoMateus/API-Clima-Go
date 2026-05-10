package provider

import (
	"api-clima-go/internal/model"
	"context"
)

type WeatherProvider interface {
	GetWeather(ctx context.Context, city string) (*model.WeatherResponse, error)
	name() string
}