package main

import (
    "fmt"
    "time"
    "github.com/vmkotov/telelog"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
    fmt.Println("🧪 Базовый тест библиотеки telelog...")
    
    // Создаем тестовое сообщение
    msg := &tgbotapi.Message{
        MessageID: 42,
        From: &tgbotapi.User{
            ID:        123456,
            FirstName: "Тест",
            UserName:  "test_user",
        },
        Chat: &tgbotapi.Chat{
            ID:    789,
            Title: "Тестовый чат",
        },
        Text: "Привет, мир!",
        Date: int(time.Now().Unix()),
    }
    
    // Создаем логгер с настройками по умолчанию
    logger := telelog.New()
    
    // Логируем сообщение
    logger.LogMessage(msg, "private")
    
    fmt.Println("\n✅ Библиотека работает!")
}
