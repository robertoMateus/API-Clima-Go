package service

import (
	"api-clima-go/internal/model"
	"api-clima-go/internal/provider"
	"context"
	"fmt"
	"time"
)

type WeatherService struct {
	providers []provider.WeatherProvider
}

func NewWeatherService(providers ...provider.WeatherProvider) * WeatherService {
	return &WeatherService{providers: providers}
}

func (s *WeatherService) GetFastest(city string) (*model.WeatherResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	resultChan := make(chan *model.WeatherResponse, len(s.providers))
	errorChan := make(chan error, len(s.providers))

	for _,p := range s.providers {
		go func (p provider.WeatherProvider) {
			result, err := p.GetWeather(ctx,city)
			if err != nil {
				errorChan <- err
				return 
			}
			resultChan <- result
		}(p)
	}

	errors := 0

	for {
		select{
		case result := <- resultChan:
			cancel()
			return result, nil

		case <-errorChan:
			errors++
			if errors == len(s.providers) {
				return nil, fmt.Errorf("todos os provedores falharam para a cidade: %s", city)
			}

			case <- ctx.Done():
				return nil, fmt.Errorf("tempo esgotado para obter dados de clima para a cidade: %s", city)
		}
	}
}