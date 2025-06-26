package weather

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// Client представляет собой клиент для работы с WeatherAPI.com
type Client struct {
	apiKey string
	client *http.Client
}

// WeatherResponse представляет ответ от API
type WeatherResponse struct {
	Location struct {
		Name    string  `json:"name"`
		Country string  `json:"country"`
		Region  string  `json:"region"`
		Lat     float64 `json:"lat"`
		Lon     float64 `json:"lon"`
	} `json:"location"`
	Current struct {
		TempC     float64 `json:"temp_c"`
		TempF     float64 `json:"temp_f"`
		Condition struct {
			Text string `json:"text"`
			Icon string `json:"icon"`
		} `json:"condition"`
		WindKph    float64 `json:"wind_kph"`
		WindDegree int     `json:"wind_degree"`
		WindDir    string  `json:"wind_dir"`
		Humidity   int     `json:"humidity"`
		Cloud      int     `json:"cloud"`
		FeelsLikeC float64 `json:"feelslike_c"`
		UV         float64 `json:"uv"`
	} `json:"current"`
}

// New создает новый клиент WeatherAPI.com
func New(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		client: &http.Client{},
	}
}

// GetWeather получает погоду для указанного города
func (c *Client) GetWeather(city string) (*WeatherResponse, error) {
	baseURL := "https://api.weatherapi.com/v1/current.json"

	// Создаем URL с параметрами
	params := url.Values{}
	params.Add("key", c.apiKey)
	params.Add("q", city)
	params.Add("aqi", "no")  // не запрашиваем данные о качестве воздуха
	params.Add("lang", "ru") // получаем ответ на русском языке

	// Формируем полный URL
	requestURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	// Отправляем запрос
	resp, err := c.client.Get(requestURL)
	if err != nil {
		return nil, fmt.Errorf("ошибка отправки запроса: %v", err)
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ошибка API: %d", resp.StatusCode)
	}

	// Декодируем ответ
	var weather WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
		return nil, fmt.Errorf("ошибка разбора ответа: %v", err)
	}

	return &weather, nil
}
