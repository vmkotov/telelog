package telelog

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// loggerImpl is the concrete implementation of TeleLogger
type loggerImpl struct {
	bot         *tgbotapi.BotAPI
	logChatID   int64
	enabled     bool
	botID       int64
	botUsername string
}

// newLoggerImpl creates a new logger instance
func newLoggerImpl(opts Options) *loggerImpl {
	return &loggerImpl{
		bot:         opts.Bot,
		logChatID:   opts.LogChatID,
		botID:       opts.BotID,
		botUsername: opts.BotUsername,
		enabled:     opts.Bot != nil && opts.LogChatID != 0,
	}
}

func (tl *loggerImpl) IsEnabled() bool {
	return tl.enabled
}

func (tl *loggerImpl) SetLogChatID(chatID int64) {
	tl.logChatID = chatID
	if chatID != 0 && tl.bot != nil {
		tl.enabled = true
	}
}

func (tl *loggerImpl) LogMessage(message *tgbotapi.Message, chatType string) {
	if !tl.enabled {
		// Fallback to console logging
		log.Printf("📨 Сообщение от @%s в %s: %s",
			message.From.UserName,
			chatType,
			message.Text)
		return
	}
	tl.sendToChat(message, "message")
}

func (tl *loggerImpl) LogCommand(message *tgbotapi.Message, command string) {
	if !tl.enabled {
		log.Printf("⚡ Команда /%s от @%s", command, message.From.UserName)
		return
	}
	tl.sendToChat(message, "command")
}

func (tl *loggerImpl) LogError(message *tgbotapi.Message, err error, context string) {
	if !tl.enabled {
		log.Printf("❌ Ошибка [%s] от @%s: %v", context, message.From.UserName, err)
		return
	}
	tl.sendErrorToChat(message, err, context)
}

func (tl *loggerImpl) SendDeployNotification(deployInfo map[string]string) {
	if !tl.enabled {
		log.Printf("🚀 Уведомление о деплое: %v", deployInfo)
		return
	}
	tl.sendDeployToChat(deployInfo)
}
