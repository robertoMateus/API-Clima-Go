package provider

import (
	"api-clima-go/internal/dto"
	"api-clima-go/internal/model"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type OpenWeatherProvider struct {
	apiKey string
	client *http.Client
}

func NewOpenWeatherProvider(apiKey string) *OpenWeatherProvider {
	return &OpenWeatherProvider{
		apiKey: apiKey,
		client: &http.Client{},
	}

}

func (p *OpenWeatherProvider) name() string {
	return "OpenWeather"
}

func (p *OpenWeatherProvider) GetWeather(ctx context.Context, city string) (*model.WeatherResponse, error) {
	url := fmt.Sprintf(
		"https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric&lang=pt_br",
		city, p.apiKey,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("Erro no OpenWeather ao criar request: %w", err)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Erro no OpenWeather na requisição: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Erro no OpenWeather: status %d", resp.StatusCode)
	}

	var data dto.OpenWeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("Erro no OpenWeather ao decodificar resposta: %w", err)
	}

	condition := ""
	if len(data.Weather) > 0 {
		condition = data.Weather[0].Description
	}

	return &model.WeatherResponse{
		City: data.Name,
		Temperature: data.Main.Temp,
		Humidity: data.Main.Humidity,
		Condition: condition,
		Source: p.name(),
	}, nil

}
