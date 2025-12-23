package telelog

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

// TeleLogger is the main logger interface
type TeleLogger interface {
	// Basic logging methods
	LogMessage(message *tgbotapi.Message, chatType string)
	LogCommand(message *tgbotapi.Message, command string)
	LogError(message *tgbotapi.Message, err error, context string)

	// Deploy notification
	SendDeployNotification(deployInfo map[string]string)

	// Configuration methods
	IsEnabled() bool
	SetLogChatID(chatID int64)
}

// Options for configuring the logger
type Options struct {
	Bot         *tgbotapi.BotAPI
	LogChatID   int64
	BotID       int64
	BotUsername string
}

// New creates a new TeleLogger instance with options
func New(opts Options) TeleLogger {
	return newLoggerImpl(opts)
}

// SimpleNew creates a simple logger without Telegram integration
// (uses console output instead)
func SimpleNew() TeleLogger {
	return newLoggerImpl(Options{})
}
