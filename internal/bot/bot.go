package bot

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"go-telegram-bot/config"
	"go-telegram-bot/internal/weather"

	tele "gopkg.in/telebot.v3"
)

// Bot представляет собой структуру нашего бота
type Bot struct {
	bot     *tele.Bot
	weather *weather.Client
	config  *config.Config
}

// WebAppData представляет данные от WebApp
type WebAppData struct {
	City      string `json:"city"`
	Timestamp int64  `json:"timestamp,omitempty"`
	UserID    int64  `json:"user_id,omitempty"`
}

// New создает новый экземпляр бота
func New(cfg *config.Config) (*Bot, error) {
	pref := tele.Settings{
		Token:  cfg.TelegramToken,
		Poller: &tele.LongPoller{Timeout: 10},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		return nil, err
	}

	w := weather.New(cfg.WeatherAPIKey)

	return &Bot{
		bot:     b,
		weather: w,
		config:  cfg,
	}, nil
}

// Start запускаем бота и HTTP сервер
func (b *Bot) Start() {
	// Запускаем HTTP сервер для обслуживания WebApp
	go b.startHTTPServer()

	// Создаем кнопку для открытия WebApp
	webapp := &tele.WebApp{
		URL: b.config.WebAppURL,
	}

	// Создаем кнопку меню для WebApp
	menuBtn := &tele.ReplyMarkup{}
	menuBtn.InlineKeyboard = [][]tele.InlineButton{
		{
			tele.InlineButton{
				Text:   "🌤 Открыть погодное приложение",
				WebApp: webapp,
			},
		},
	}

	// Обработчик команды /start
	b.bot.Handle("/start", func(c tele.Context) error {
		message := "Привет! Я учебный бот на Go.\n\n" +
			"Используйте /help для получения списка команд или нажмите кнопку ниже для открытия погодного приложения:"
		return c.Send(message, menuBtn)
	})

	// Обработчик команды /help
	b.bot.Handle("/help", func(c tele.Context) error {
		helpText := `Доступные команды:
/start - Начать работу с ботом
/help - Показать это сообщение
/echo <текст> - Повторить ваш текст
/weather <город> - Показать погоду в городе
/app - Открыть погодное приложение`
		return c.Send(helpText)
	})

	// Обработчик команды /app для открытия WebApp
	b.bot.Handle("/app", func(c tele.Context) error {
		message := "Нажмите кнопку ниже для открытия погодного приложения:"
		return c.Send(message, menuBtn)
	})

	// Обработчик команды /echo
	b.bot.Handle("/echo", func(c tele.Context) error {
		text := c.Message().Payload
		if text == "" {
			return c.Send("Пожалуйста, укажите текст после команды /echo")
		}
		return c.Send(text)
	})

	// Обработчик команды /weather
	b.bot.Handle("/weather", func(c tele.Context) error {
		city := c.Message().Payload
		if city == "" {
			return c.Send("Пожалуйста, укажите город после команды /weather\nНапример: /weather Москва")
		}

		return b.sendWeatherResponse(c, city)
	})

	// Обработчик данных от WebApp
	b.bot.Handle(tele.OnWebApp, func(c tele.Context) error {
		var data WebAppData
		if err := json.Unmarshal([]byte(c.Message().WebAppData.Data), &data); err != nil {
			log.Printf("Ошибка разбора данных WebApp: %v", err)
			return c.Send("❌ Ошибка при обработке данных приложения")
		}

		if strings.TrimSpace(data.City) == "" {
			return c.Send("❌ Пожалуйста, укажите город")
		}

		return b.sendWeatherResponse(c, data.City)
	})

	log.Println("Bot started...")
	b.bot.Start()
}

// sendWeatherResponse отправляет ответ с информацией о погоде
func (b *Bot) sendWeatherResponse(c tele.Context, city string) error {
	// Получаем погоду
	weather, err := b.weather.GetWeather(city)
	if err != nil {
		return c.Send(fmt.Sprintf("❌ Ошибка при получении погоды: %v", err))
	}

	// Формируем сообщение с погодой
	message := fmt.Sprintf("🌡 Погода в %s, %s:\n"+
		"Температура: %.1f°C (ощущается как %.1f°C)\n"+
		"Погодные условия: %s\n"+
		"Влажность: %d%%\n"+
		"Ветер: %.1f км/ч, %s\n"+
		"Облачность: %d%%\n"+
		"УФ-индекс: %.1f",
		weather.Location.Name,
		weather.Location.Country,
		weather.Current.TempC,
		weather.Current.FeelsLikeC,
		weather.Current.Condition.Text,
		weather.Current.Humidity,
		weather.Current.WindKph,
		weather.Current.WindDir,
		weather.Current.Cloud,
		weather.Current.UV,
	)

	return c.Send(message)
}

// startHTTPServer запускает HTTP сервер для обслуживания WebApp
func (b *Bot) startHTTPServer() {
	// Обслуживание статических файлов
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static/"))))

	// Обработчик для главной страницы WebApp
	http.HandleFunc("/", b.serveWebApp)

	// Обработчик для API endpoint (если нужен)
	http.HandleFunc("/api/weather", b.handleWeatherAPI)

	log.Printf("Starting HTTP server on :%s...", b.config.HTTPPort)
	if err := http.ListenAndServe(":"+b.config.HTTPPort, nil); err != nil {
		log.Printf("HTTP server error: %v", err)
	}
}

// serveWebApp обслуживает HTML страницу WebApp
func (b *Bot) serveWebApp(w http.ResponseWriter, r *http.Request) {
	// Читаем HTML файл
	htmlContent, err := http.Dir("web").Open("index.html")
	if err != nil {
		log.Printf("Error opening HTML file: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer htmlContent.Close()

	// Устанавливаем заголовки
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	// Копируем содержимое файла в ответ
	http.ServeContent(w, r, "index.html", time.Time{}, htmlContent)
}

// handleWeatherAPI обрабатывает API запросы (если нужно)
func (b *Bot) handleWeatherAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "ok", "message": "Weather API endpoint"}`))
}
