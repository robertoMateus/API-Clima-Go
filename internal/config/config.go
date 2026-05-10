package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	OpenWeatherApiKey string
	WeatherApiKey     string
	ServerPort        string
}

func getEnv(key string) string {
	value := os.Getenv(key)

	if value == "" {
		log.Fatalf("Variável de ambiente obrigatória não encontrada: %s", key)
	}
	return value
}

func Load() *Config{
	if err := godotenv.Load(); err != nil {
		log.Println("Arquivo .env não encontrado!")
	}

	return &Config{
		OpenWeatherApiKey: getEnv("OPEN_WEATHER_API_KEY"),
		WeatherApiKey:     getEnv("WEATHER_API_KEY"),
		ServerPort:        getEnv("SERVER_PORT"),
	}
}