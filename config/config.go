package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config содержит все конфигурационные параметры
type Config struct {
	TelegramToken string
	WeatherAPIKey string
	WebAppURL     string
	HTTPPort      string
}

// Load загружает конфигурацию из переменных окружения
func Load() *Config {
	// Загружаем .env файл
	if err := godotenv.Load(); err != nil {
		// Если файл не найден, это не ошибка - просто используем переменные окружения
		// которые уже установлены в системе
	}

	webAppURL := os.Getenv("WEBAPP_URL")
	if webAppURL == "" {
		webAppURL = "http://localhost:8080" // Значение по умолчанию для разработки
	}

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080" // Порт по умолчанию
	}

	return &Config{
		TelegramToken: os.Getenv("TELEGRAM_TOKEN"),
		WeatherAPIKey: os.Getenv("WEATHER_API_KEY"),
		WebAppURL:     webAppURL,
		HTTPPort:      httpPort,
	}
}
