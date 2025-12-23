package telelog

import (
	"fmt"
	"log"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// TeleLogger представляет логгер для Telegram
type TeleLogger struct {
	bot         *tgbotapi.BotAPI
	logChatID   int64
	enabled     bool
	botID       int64  // ID вашего бота
	botUsername string // Username вашего бота
}

// New создает новый логгер
func New(bot *tgbotapi.BotAPI, logChatID, botID int64, botUsername string) *TeleLogger {
	return &TeleLogger{
		bot:         bot,
		logChatID:   logChatID,
		botID:       botID,
		botUsername: botUsername,
		enabled:     bot != nil && logChatID != 0,
	}
}

// SimpleNew создает упрощенный логгер (для обратной совместимости)
func SimpleNew() *TeleLogger {
	return &TeleLogger{
		enabled: false,
	}
}

// SendToChat отправляет сообщение в указанный чат для логирования
func (tl *TeleLogger) SendToChat(message *tgbotapi.Message, targetChatID int64) {
	if !tl.enabled || tl.bot == nil {
		return
	}

	// Проверяем, что это не сообщение от самого бота
	if message.From != nil && message.From.ID == tl.botID {
		return
	}

	// Формируем информативное сообщение
	chatInfo := tl.formatChatInfo(message)
	userInfo := tl.formatUserInfo(message)
	messageInfo := tl.formatMessageInfo(message)
	botInfo := tl.formatBotInfo()

	// Собираем итоговое сообщение
	text := fmt.Sprintf(
		"🤖 *Лог сообщения* %s\n\n"+
			"%s\n"+
			"%s\n"+
			"%s\n"+
			"%s",
		message.Time().Format("15:04:05"),
		chatInfo,
		userInfo,
		messageInfo,
		botInfo,
	)

	// Ограничиваем длину сообщения (Telegram лимит ~4096 символов)
	if len(text) > 4000 {
		text = text[:4000] + "\n... (сообщение обрезано)"
	}

	msg := tgbotapi.NewMessage(targetChatID, text)
	msg.ParseMode = "Markdown"

	if _, err := tl.bot.Send(msg); err != nil {
		log.Printf("❌ Не удалось отправить логи в чат %d: %v", targetChatID, err)
	} else {
		log.Printf("✅ Логи отправлены в чат %d", targetChatID)
	}
}

// LogMessage логирует сообщение (старый метод для совместимости)
func (tl *TeleLogger) LogMessage(message *tgbotapi.Message, chatType string) {
	if !tl.enabled {
		log.Printf("📨 Сообщение от @%s в %s: %s",
			message.From.UserName,
			chatType,
			message.Text)
		return
	}

	// Отправляем в лог-чат если есть бот
	tl.SendToChat(message, tl.logChatID)
}

// LogCommand логирует команду (старый метод для совместимости)
func (tl *TeleLogger) LogCommand(message *tgbotapi.Message, command string) {
	if !tl.enabled {
		log.Printf("⚡ Команда /%s от @%s", command, message.From.UserName)
		return
	}

	// Отправляем в лог-чат если есть бот
	tl.SendToChat(message, tl.logChatID)
}

// formatChatInfo форматирует информацию о чате
func (tl *TeleLogger) formatChatInfo(message *tgbotapi.Message) string {
	chatType := "личный"
	if message.Chat.IsGroup() {
		chatType = "группа"
	} else if message.Chat.IsSuperGroup() {
		chatType = "супергруппа"
	} else if message.Chat.IsChannel() {
		chatType = "канал"
	}

	chatTitle := "Без названия"
	if message.Chat.Title != "" {
		chatTitle = message.Chat.Title
	}

	return fmt.Sprintf(
		"💬 *Чат:* %s\n"+
			"📌 Тип: %s\n"+
			"🆔 ID: `%d`",
		chatTitle,
		chatType,
		message.Chat.ID,
	)
}

// formatUserInfo форматирует информацию о пользователе
func (tl *TeleLogger) formatUserInfo(message *tgbotapi.Message) string {
	if message.From == nil {
		return "👤 *Пользователь:* Неизвестен"
	}

	userName := message.From.UserName
	if userName == "" {
		userName = "без username"
	}

	fullName := fmt.Sprintf("%s %s",
		message.From.FirstName,
		message.From.LastName)
	fullName = strings.TrimSpace(fullName)
	if fullName == "" {
		fullName = "Без имени"
	}

	return fmt.Sprintf(
		"👤 *Пользователь:* %s\n"+
			"📛 Имя: %s\n"+
			"🔖 @%s\n"+
			"🆔 ID: `%d`",
		fullName,
		message.From.FirstName,
		userName,
		message.From.ID,
	)
}

// formatMessageInfo форматирует информацию о сообщении
func (tl *TeleLogger) formatMessageInfo(message *tgbotapi.Message) string {
	messageText := message.Text
	if messageText == "" {
		messageText = "⚠️ *Без текста*"

		// Проверяем другие типы контента
		if message.Sticker != nil {
			messageText = fmt.Sprintf("🎭 Стикер: %s", message.Sticker.Emoji)
		} else if message.Photo != nil && len(message.Photo) > 0 {
			messageText = "🖼️ Фото"
		} else if message.Video != nil {
			messageText = "🎬 Видео"
		} else if message.Document != nil {
			messageText = fmt.Sprintf("📎 Документ: %s", message.Document.FileName)
		} else if message.Audio != nil {
			messageText = "🎵 Аудио"
		} else if message.Voice != nil {
			messageText = "🎤 Голосовое сообщение"
		} else if message.Location != nil {
			messageText = "📍 Локация"
		} else if message.Contact != nil {
			messageText = "👤 Контакт"
		}
	}

	info := fmt.Sprintf("📝 *Сообщение:*\n%s", messageText)

	// Добавляем информацию о reply, если есть
	if message.ReplyToMessage != nil {
		replyText := message.ReplyToMessage.Text
		if replyText == "" {
			replyText = "⬆️ (сообщение без текста)"
		}
		if len(replyText) > 100 {
			replyText = replyText[:100] + "..."
		}

		info += fmt.Sprintf("\n\n↩️ *Ответ на:*\n%s", replyText)
	}

	return info
}

// formatBotInfo форматирует информацию о боте
func (tl *TeleLogger) formatBotInfo() string {
	return fmt.Sprintf(
		"\n🤖 *Информация о боте:*\n"+
			"Бот: @%s\n"+
			"Bot ID: `%d`\n"+
			"⏰ Лог создан: %s",
		tl.botUsername,
		tl.botID,
		time.Now().Format("2006-01-02 15:04:05"),
	)
}

// IsEnabled возвращает статус логгера
func (tl *TeleLogger) IsEnabled() bool {
	return tl.enabled
}

// SetLogChatID устанавливает ID чата для логов
func (tl *TeleLogger) SetLogChatID(chatID int64) {
	tl.logChatID = chatID
	if chatID != 0 && tl.bot != nil {
		tl.enabled = true
	}
}
