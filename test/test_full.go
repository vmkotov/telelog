package main

import (
    "fmt"
    "time"
    "github.com/vmkotov/telelog"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
    fmt.Println("🧪 Расширенный тест библиотеки telelog...")
    
    // 1. Сообщение с командой
    msg1 := &tgbotapi.Message{
        MessageID: 1,
        From: &tgbotapi.User{
            ID:        266468924,
            FirstName: "Вячеслав",
            UserName:  "vmkotov",
        },
        Chat: &tgbotapi.Chat{
            ID:    -1001234567890,
            Title: "Bushlatinga Chat",
        },
        Text: "/start",
        Date: int(time.Now().Unix()),
    }
    
    // 2. Сообщение со стикером
    msg2 := &tgbotapi.Message{
        MessageID: 2,
        From: &tgbotapi.User{
            ID:        555555555,
            FirstName: "Другой",
            UserName:  "friend",
        },
        Chat: &tgbotapi.Chat{
            ID:    -1001234567890,
            Title: "Bushlatinga Chat",
        },
        Date: int(time.Now().Unix()),
        Sticker: &tgbotapi.Sticker{
            FileID:       "CAACAgIAAxkBAAN",
            FileUniqueID: "test_sticker",
            Emoji:        "😎",
        },
    }
    
    logger := telelog.New()
    
    fmt.Println("\n=== Тест 1: Обычное сообщение с командой ===")
    logger.LogMessage(msg1, "supergroup")
    logger.LogCommand(msg1, "start")
    
    fmt.Println("\n=== Тест 2: Сообщение со стикером ===")
    logger.LogMessage(msg2, "supergroup")
    
    fmt.Println("\n=== Тест 3: Логирование ошибки ===")
    logger.LogError(msg1, fmt.Errorf("не удалось подключиться к БД"), "обработка команды /start")
    
    fmt.Println("\n=== Тест 4: Компактный формат (без эмодзи) ===")
    compactLogger := telelog.New(telelog.Options{
        Format:   "compact",
        Colorize: false,
    })
    msg1.Text = "Славик привет!"
    compactLogger.LogMessage(msg1, "supergroup")
    
    fmt.Println("\n✅ Все тесты пройдены!")
}
