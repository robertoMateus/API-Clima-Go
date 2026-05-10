package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"api-clima-go/internal/dto"
	"api-clima-go/internal/model"
)

type WeatherAPIProvider struct {
	apiKey string
	client *http.Client
}

func NewWeatherAPIProvider(apiKey string) *WeatherAPIProvider {
	return &WeatherAPIProvider{
		apiKey: apiKey,
		client: &http.Client{},
	}
}

func (p *WeatherAPIProvider) name() string {
	return "WeatherAPI"
}

func (p *WeatherAPIProvider) GetWeather(ctx context.Context, city string) (*model.WeatherResponse, error) {
	url := fmt.Sprintf(
		"https://api.weatherapi.com/v1/current.json?key=%s&q=%s&lang=pt",
		p.apiKey, city,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("erro no WeatherAPI ao criar request: %w", err)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro no WeatherAPI na requisição: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("erro no WeatherAPI: status inesperado: %d", resp.StatusCode)
	}

	var data dto.WeatherApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("erro no WeatherAPI ao decodificar resposta: %w", err)
	}

	return &model.WeatherResponse{
		City:        data.Location.Name,
		Temperature: data.Current.TempC,
		Humidity:    data.Current.Humidity,
		Condition:   data.Current.Condition.Text,
		Source:      p.name(),
	}, nil
}