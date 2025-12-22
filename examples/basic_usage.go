package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vmkotov/telelog"
)

func main() {
	// Создаем логгер с опциями по умолчанию
	logger := telelog.New()

	// Или с кастомными опциями
	customLogger := telelog.New(telelog.Options{
		Format:        "compact",
		IncludeEmojis: false,
		Colorize:      false,
	})

	// Или с кастомным writer
	fileLogger := telelog.New(telelog.Options{
		Writer: log.New(os.Stdout, "[BOT] ", log.LstdFlags),
	})

	// Имитация входящего сообщения
	msg := &tgbotapi.Message{
		MessageID: 123,
		From: &tgbotapi.User{
			ID:           456,
			FirstName:    "John",
			LastName:     "Doe",
			UserName:     "johndoe",
			LanguageCode: "en",
			IsBot:        false,
		},
		Chat: &tgbotapi.Chat{
			ID:    789,
			Title: "Test Chat",
			Type:  "private",
		},
		Text: "Hello, world!",
		Date: 1672531200, // Unix timestamp
	}

	// Использование
	logger.LogMessage(msg, "private")
	customLogger.LogMessage(msg, "private")
	fileLogger.LogMessage(msg, "private")

	// Логирование команды
	logger.LogCommand(msg, "start")

	// Логирование ошибки
	logger.LogError(msg, nil, "test error")
}
