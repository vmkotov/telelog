package main

import (
	"fmt"
	"time"
	
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vmkotov/telelog"
)

func main() {
	fmt.Println("🧪 Пример использования библиотеки telelog")
	
	// 1. Простой логгер (консольный)
	fmt.Println("\n1. Простой консольный логгер:")
	consoleLogger := telelog.SimpleNew()
	
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
		Date: int(time.Now().Unix()),
	}
	
	consoleLogger.LogMessage(msg, "private")
	consoleLogger.LogCommand(msg, "start")
	consoleLogger.LogError(msg, fmt.Errorf("example error"), "test operation")
	
	// 2. Уведомление о деплое
	fmt.Println("\n2. Уведомление о деплое:")
	deployInfo := map[string]string{
		"version":     "1.0.0",
		"environment": "staging",
		"branch":      "feature/telelog",
		"commit_hash": "abc123def456",
		"deployer":    "CI/CD Pipeline",
		"timestamp":   time.Now().Format("2006-01-02 15:04:05"),
	}
	consoleLogger.SendDeployNotification(deployInfo)
	
	// 3. Логгер с Telegram ботом (пример, нужен реальный токен)
	fmt.Println("\n3. Пример с Telegram ботом (закомментирован):")
	/*
	bot, err := tgbotapi.NewBotAPI("YOUR_BOT_TOKEN")
	if err != nil {
		log.Fatal(err)
	}
	
	telegramLogger := telelog.New(telelog.Options{
		Bot:         bot,
		LogChatID:   -1001234567890, // ID чата для логов
		BotID:       bot.Self.ID,
		BotUsername: bot.Self.UserName,
	})
	
	telegramLogger.LogMessage(msg, "private")
	*/
	
	fmt.Println("\n✅ Пример завершён!")
}
