package telelog

import (
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (tl *loggerImpl) sendToChat(message *tgbotapi.Message, messageType string) {
	if !tl.enabled || tl.bot == nil {
		return
	}

	if message.From != nil && message.From.ID == tl.botID {
		return
	}

	chatInfo := tl.formatChatInfo(message)
	userInfo := tl.formatUserInfo(message)
	messageInfo := tl.formatMessageInfo(message)
	botInfo := tl.formatBotInfo()

	// Добавляем тип сообщения в заголовок
	title := "🤖 Лог сообщения"
	if messageType == "command" {
		title = "⚡ Лог команды"
	}

	text := fmt.Sprintf(
		"%s %s\n\n"+
			"%s\n"+
			"%s\n"+
			"%s\n"+
			"%s",
		title,
		message.Time().Format("15:04:05"),
		chatInfo,
		userInfo,
		messageInfo,
		botInfo,
	)

	if len(text) > 4000 {
		text = text[:4000] + "\n... (сообщение обрезано)"
	}

	msg := tgbotapi.NewMessage(tl.logChatID, text)

	if _, err := tl.bot.Send(msg); err != nil {
		log.Printf("❌ Не удалось отправить логи в чат %d: %v", tl.logChatID, err)
	} else {
		log.Printf("✅ Логи отправлены в чат %d", tl.logChatID)
	}
}

func (tl *loggerImpl) sendErrorToChat(message *tgbotapi.Message, err error, context string) {
	if !tl.enabled || tl.bot == nil {
		return
	}

	errorMsg := "Неизвестная ошибка"
	if err != nil {
		errorMsg = err.Error()
	}

	text := fmt.Sprintf(
		"🚨 *ОШИБКА*\n\n"+
			"*Контекст:* %s\n"+
			"*Сообщение:* %s\n"+
			"*Пользователь:* @%s (%s)\n"+
			"*Чат:* %s\n"+
			"*Время:* %s\n\n"+
			"*Ошибка:* `%s`",
		context,
		message.Text,
		message.From.UserName,
		message.From.FirstName,
		tl.formatChatTitle(message),
		time.Now().Format("15:04:05"),
		errorMsg,
	)

	msg := tgbotapi.NewMessage(tl.logChatID, text)
	msg.ParseMode = "Markdown"

	if _, err := tl.bot.Send(msg); err != nil {
		log.Printf("❌ Не удалось отправить ошибку в чат %d: %v", tl.logChatID, err)
	}
}

// Helper methods
func (tl *loggerImpl) formatChatInfo(message *tgbotapi.Message) string {
	chatTitle := "Неизвестный чат"
	if message.Chat != nil && message.Chat.Title != "" {
		chatTitle = message.Chat.Title
	}
	return fmt.Sprintf("💬 *Чат:* %s (ID: %d)", chatTitle, message.Chat.ID)
}

func (tl *loggerImpl) formatUserInfo(message *tgbotapi.Message) string {
	if message.From == nil {
		return "👤 *Пользователь:* Неизвестен"
	}
	
	lastName := ""
	if message.From.LastName != "" {
		lastName = " " + message.From.LastName
	}
	
	userName := ""
	if message.From.UserName != "" {
		userName = fmt.Sprintf(" (@%s)", message.From.UserName)
	}
	
	return fmt.Sprintf("👤 *Пользователь:* %s%s%s (ID: %d)",
		message.From.FirstName,
		lastName,
		userName,
		message.From.ID)
}

func (tl *loggerImpl) formatMessageInfo(message *tgbotapi.Message) string {
	if message.Sticker != nil {
		return fmt.Sprintf("🎭 *Стикер:* %s", message.Sticker.Emoji)
	}
	if message.Text != "" {
		return fmt.Sprintf("💭 *Текст:* %s", message.Text)
	}
	if message.Photo != nil {
		return "🖼️ *Фото*"
	}
	if message.Document != nil {
		return fmt.Sprintf("📎 *Документ:* %s", message.Document.FileName)
	}
	return "📦 *Сообщение без текста*"
}

func (tl *loggerImpl) formatBotInfo() string {
	botUsername := "unknown"
	if tl.botUsername != "" {
		botUsername = tl.botUsername
	}
	return fmt.Sprintf("🤖 *Бот:* @%s", botUsername)
}

func (tl *loggerImpl) formatChatTitle(message *tgbotapi.Message) string {
	if message.Chat != nil && message.Chat.Title != "" {
		return message.Chat.Title
	}
	if message.Chat != nil && message.Chat.UserName != "" {
		return "@" + message.Chat.UserName
	}
	return fmt.Sprintf("ID: %d", message.Chat.ID)
}
