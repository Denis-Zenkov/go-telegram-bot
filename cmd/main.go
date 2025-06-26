package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"go-telegram-bot/config"
	"go-telegram-bot/internal/bot"
)

func main() {
	// Загружаем конфигурацию
	cfg := config.Load()
	if cfg.TelegramToken == "" {
		log.Fatal("TELEGRAM_TOKEN не установлен")
	}

	// Создаем и запускаем бота
	b, err := bot.New(cfg)
	if err != nil {
		log.Fatalf("Ошибка создания бота: %v", err)
	}

	// Создаем канал для получения сигналов операционной системы
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	log.Println("Starting bot...")

	// Запускаем бота в отдельной горутине
	go b.Start()

	// Ожидаем сигнал завершения
	<-sigChan
	log.Println("Shutting down...")
}
