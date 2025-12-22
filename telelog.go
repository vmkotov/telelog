package telelog

import (
	"fmt"
	"log"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Options - опции для логирования
type Options struct {
	// Форматы вывода
	Format        string // "text", "json", "compact"
	IncludeEmojis bool
	IncludeDate   bool
	IncludeUser   bool
	IncludeChat   bool
	IncludeMedia  bool
	Colorize      bool
	LogLevel      string // "info", "debug", "error"
	
	// Кастомные функции
	CustomFormatter func(*MessageDetails) string
	Writer          *log.Logger
}

// MessageDetails - детали сообщения
type MessageDetails struct {
	Message    *tgbotapi.Message
	ChatType   string
	Timestamp  time.Time
	MessageID  int
	FromUser   UserInfo
	ChatInfo   ChatInfo
	MediaInfo  *MediaInfo
}

// UserInfo - информация о пользователе
type UserInfo struct {
	ID           int64
	FirstName    string
	LastName     string
	Username     string
	LanguageCode string
	IsBot        bool
}

// ChatInfo - информация о чате
type ChatInfo struct {
	ID    int64
	Title string
	Type  string // "private", "group", "supergroup", "channel"
}

// MediaInfo - информация о медиа
type MediaInfo struct {
	Type      string // "photo", "sticker", "document", "voice", "location", "video"
	FileID    string
	FileName  string
	MimeType  string
	Emoji     string
	Duration  int
	Latitude  float64
	Longitude float64
}

// Logger - основной логгер
type Logger struct {
	options Options
}

// DefaultOptions - опции по умолчанию
var DefaultOptions = Options{
	Format:        "text",
	IncludeEmojis: true,
	IncludeDate:   true,
	IncludeUser:   true,
	IncludeChat:   true,
	IncludeMedia:  true,
	Colorize:      true,
	LogLevel:      "info",
}

// New создает новый логгер
func New(options ...Options) *Logger {
	if len(options) > 0 {
		return &Logger{options: options[0]}
	}
	return &Logger{options: DefaultOptions}
}

// LogMessage логирует детали сообщения
func (l *Logger) LogMessage(message *tgbotapi.Message, chatType string) {
	details := l.extractDetails(message, chatType)
	
	var output string
	
	if l.options.CustomFormatter != nil {
		output = l.options.CustomFormatter(details)
	} else {
		switch l.options.Format {
		case "json":
			output = l.formatJSON(details)
		case "compact":
			output = l.formatCompact(details)
		default: // "text"
			output = l.formatText(details)
		}
	}
	
	writer := l.options.Writer
	if writer == nil {
		writer = log.Default()
	}
	
	writer.Println(output)
}

// LogCommand логирует команду отдельно
func (l *Logger) LogCommand(message *tgbotapi.Message, command string) {
	details := l.extractDetails(message, getChatType(message.Chat))
	output := l.formatCommand(details, command)
	
	writer := l.options.Writer
	if writer == nil {
		writer = log.Default()
	}
	
	writer.Println(output)
}

// LogError логирует ошибку с контекстом сообщения
func (l *Logger) LogError(message *tgbotapi.Message, err error, context string) {
	details := l.extractDetails(message, getChatType(message.Chat))
	output := l.formatError(details, err, context)
	
	writer := l.options.Writer
	if writer == nil {
		writer = log.Default()
	}
	
	writer.Println(output)
}

// extractDetails извлекает детали из сообщения
func (l *Logger) extractDetails(message *tgbotapi.Message, chatType string) *MessageDetails {
	details := &MessageDetails{
		Message:   message,
		ChatType:  chatType,
		Timestamp: time.Unix(int64(message.Date), 0),
		MessageID: message.MessageID,
		FromUser: UserInfo{
			ID:           message.From.ID,
			FirstName:    message.From.FirstName,
			LastName:     message.From.LastName,
			Username:     message.From.UserName,
			LanguageCode: message.From.LanguageCode,
			IsBot:        message.From.IsBot,
		},
		ChatInfo: ChatInfo{
			ID:    message.Chat.ID,
			Title: getChatTitle(message.Chat),
			Type:  chatType,
		},
	}
	
	// Извлекаем медиа информацию
	if message.Photo != nil && len(message.Photo) > 0 {
		details.MediaInfo = &MediaInfo{
			Type:   "photo",
			FileID: message.Photo[len(message.Photo)-1].FileID,
		}
	} else if message.Sticker != nil {
		details.MediaInfo = &MediaInfo{
			Type:   "sticker",
			FileID: message.Sticker.FileID,
			Emoji:  message.Sticker.Emoji,
		}
	} else if message.Document != nil {
		details.MediaInfo = &MediaInfo{
			Type:     "document",
			FileID:   message.Document.FileID,
			FileName: message.Document.FileName,
			MimeType: message.Document.MimeType,
		}
	} else if message.Voice != nil {
		details.MediaInfo = &MediaInfo{
			Type:     "voice",
			FileID:   message.Voice.FileID,
			Duration: message.Voice.Duration,
		}
	} else if message.Location != nil {
		details.MediaInfo = &MediaInfo{
			Type:      "location",
			Latitude:  message.Location.Latitude,
			Longitude: message.Location.Longitude,
		}
	} else if message.Video != nil {
		details.MediaInfo = &MediaInfo{
			Type:     "video",
			FileID:   message.Video.FileID,
			Duration: message.Video.Duration,
		}
	}
	
	return details
}

// formatText форматирует в читаемый текст
func (l *Logger) formatText(details *MessageDetails) string {
	var builder strings.Builder
	
	if l.options.Colorize {
		builder.WriteString("\033[36m📥\033[0m ") // Голубой emoji
	} else if l.options.IncludeEmojis {
		builder.WriteString("📥 ")
	}
	
	builder.WriteString("INCOMING MESSAGE:\n")
	
	if l.options.IncludeDate {
		if l.options.Colorize {
			builder.WriteString(fmt.Sprintf("   \033[33m📅\033[0m Date: %s\n", details.Timestamp.Format("2006-01-02 15:04:05")))
		} else if l.options.IncludeEmojis {
			builder.WriteString(fmt.Sprintf("   📅 Date: %s\n", details.Timestamp.Format("2006-01-02 15:04:05")))
		} else {
			builder.WriteString(fmt.Sprintf("   Date: %s\n", details.Timestamp.Format("2006-01-02 15:04:05")))
		}
	}
	
	if l.options.IncludeUser {
		if l.options.Colorize {
			builder.WriteString(fmt.Sprintf("   \033[34m👤\033[0m User: %s %s (ID: %d, @%s)\n", 
				details.FromUser.FirstName, 
				details.FromUser.LastName,
				details.FromUser.ID,
				details.FromUser.Username))
		} else if l.options.IncludeEmojis {
			builder.WriteString(fmt.Sprintf("   👤 User: %s %s (ID: %d, @%s)\n", 
				details.FromUser.FirstName, 
				details.FromUser.LastName,
				details.FromUser.ID,
				details.FromUser.Username))
		} else {
			builder.WriteString(fmt.Sprintf("   User: %s %s (ID: %d, @%s)\n", 
				details.FromUser.FirstName, 
				details.FromUser.LastName,
				details.FromUser.ID,
				details.FromUser.Username))
		}
	}
	
	if l.options.IncludeChat {
		if l.options.Colorize {
			builder.WriteString(fmt.Sprintf("   \033[35m💬\033[0m Chat: %s (ID: %d, Type: %s)\n", 
				details.ChatInfo.Title,
				details.ChatInfo.ID,
				details.ChatInfo.Type))
		} else if l.options.IncludeEmojis {
			builder.WriteString(fmt.Sprintf("   💬 Chat: %s (ID: %d, Type: %s)\n", 
				details.ChatInfo.Title,
				details.ChatInfo.ID,
				details.ChatInfo.Type))
		} else {
			builder.WriteString(fmt.Sprintf("   Chat: %s (ID: %d, Type: %s)\n", 
				details.ChatInfo.Title,
				details.ChatInfo.ID,
				details.ChatInfo.Type))
		}
	}
	
	// Текст сообщения
	if details.Message.Text != "" {
		if l.options.Colorize {
			builder.WriteString(fmt.Sprintf("   \033[32m📝\033[0m Text: %s\n", details.Message.Text))
		} else if l.options.IncludeEmojis {
			builder.WriteString(fmt.Sprintf("   📝 Text: %s\n", details.Message.Text))
		} else {
			builder.WriteString(fmt.Sprintf("   Text: %s\n", details.Message.Text))
		}
	}
	
	if l.options.IncludeMedia && details.MediaInfo != nil {
		switch details.MediaInfo.Type {
		case "photo":
			builder.WriteString("   📸 Photo\n")
		case "sticker":
			builder.WriteString(fmt.Sprintf("   🎭 Sticker: %s\n", details.MediaInfo.Emoji))
		case "document":
			builder.WriteString(fmt.Sprintf("   📎 Document: %s\n", details.MediaInfo.FileName))
		case "voice":
			builder.WriteString(fmt.Sprintf("   🎤 Voice: %d sec\n", details.MediaInfo.Duration))
		case "location":
			builder.WriteString(fmt.Sprintf("   📍 Location: lat=%.6f, lon=%.6f\n", 
				details.MediaInfo.Latitude, details.MediaInfo.Longitude))
		case "video":
			builder.WriteString("   🎬 Video\n")
		}
	}
	
	// Дополнительная информация
	if details.Message.ReplyToMessage != nil {
		builder.WriteString(fmt.Sprintf("   ↪️  Reply to: %d\n", details.Message.ReplyToMessage.MessageID))
	}
	
	if details.Message.ForwardFrom != nil {
		builder.WriteString(fmt.Sprintf("   ↩️  Forwarded from user ID: %d\n", details.Message.ForwardFrom.ID))
	}
	
	return strings.TrimSuffix(builder.String(), "\n")
}

// SimpleFormatter - простой форматтер без лишней информации
func (l *Logger) SimpleFormatter(details *MessageDetails) string {
	return fmt.Sprintf("%s: %s", details.FromUser.Username, details.Message.Text)
}

// DebugFormatter - детальный форматтер для отладки
func (l *Logger) DebugFormatter(details *MessageDetails) string {
	var builder strings.Builder
	
	builder.WriteString("=== DEBUG MESSAGE ===\n")
	builder.WriteString(fmt.Sprintf("User: ID=%d, First=%s, Last=%s, @%s, Lang=%s, Bot=%v\n",
		details.FromUser.ID,
		details.FromUser.FirstName,
		details.FromUser.LastName,
		details.FromUser.Username,
		details.FromUser.LanguageCode,
		details.FromUser.IsBot))
	
	builder.WriteString(fmt.Sprintf("Chat: ID=%d, Title=%s, Type=%s\n",
		details.ChatInfo.ID,
		details.ChatInfo.Title,
		details.ChatInfo.Type))
	
	builder.WriteString(fmt.Sprintf("Time: %s\n", details.Timestamp.Format("2006-01-02 15:04:05.000")))
	builder.WriteString(fmt.Sprintf("Text: %s\n", details.Message.Text))
	
	if details.MediaInfo != nil {
		builder.WriteString(fmt.Sprintf("Media: Type=%s\n", details.MediaInfo.Type))
	}
	
	if details.Message.ReplyToMessage != nil {
		builder.WriteString(fmt.Sprintf("ReplyTo: %d\n", details.Message.ReplyToMessage.MessageID))
	}
	
	return builder.String()
}

// formatCompact компактный формат
func (l *Logger) formatCompact(details *MessageDetails) string {
	timestamp := details.Timestamp.Format("15:04:05")
	user := fmt.Sprintf("@%s", details.FromUser.Username)
	if user == "@" {
		user = fmt.Sprintf("%d", details.FromUser.ID)
	}
	
	return fmt.Sprintf("[%s] %s (%s): %s", 
		timestamp, 
		user, 
		details.ChatInfo.Type, 
		truncate(details.Message.Text, 50))
}

// formatJSON JSON формат (упрощенный)
func (l *Logger) formatJSON(details *MessageDetails) string {
	text := escapeJSON(details.Message.Text)
	username := details.FromUser.Username
	if username == "" {
		username = fmt.Sprintf("%d", details.FromUser.ID)
	}
	
	return fmt.Sprintf(`{"timestamp":"%s","user_id":%d,"username":"%s","chat_id":%d,"text":"%s"}`,
		details.Timestamp.Format(time.RFC3339),
		details.FromUser.ID,
		username,
		details.ChatInfo.ID,
		text)
}

// formatCommand форматирует команду
func (l *Logger) formatCommand(details *MessageDetails, command string) string {
	if l.options.Colorize {
		return fmt.Sprintf("\033[35m⚡\033[0m Command from @%s in %s: /%s", 
			details.FromUser.Username, 
			details.ChatInfo.Title, 
			command)
	}
	return fmt.Sprintf("⚡ Command from @%s in %s: /%s", 
		details.FromUser.Username, 
		details.ChatInfo.Title, 
		command)
}

// formatError форматирует ошибку
func (l *Logger) formatError(details *MessageDetails, err error, context string) string {
	if l.options.Colorize {
		return fmt.Sprintf("\033[31m❌\033[0m Error in %s: %v (Context: %s)", 
			details.ChatInfo.Title, 
			err, 
			context)
	}
	return fmt.Sprintf("❌ Error in %s: %v (Context: %s)", 
		details.ChatInfo.Title, 
		err, 
		context)
}

// Helper functions
func getChatTitle(chat *tgbotapi.Chat) string {
	if chat.Title != "" {
		return chat.Title
	}
	if chat.FirstName != "" {
		title := chat.FirstName
		if chat.LastName != "" {
			title += " " + chat.LastName
		}
		return title
	}
	return "Unknown"
}

func getChatType(chat *tgbotapi.Chat) string {
	if chat.IsPrivate() {
		return "private"
	}
	if chat.IsGroup() {
		return "group"
	}
	if chat.IsSuperGroup() {
		return "supergroup"
	}
	if chat.IsChannel() {
		return "channel"
	}
	return "unknown"
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func escapeJSON(s string) string {
	s = strings.ReplaceAll(s, `"`, `\"`)
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, "\n", `\n`)
	s = strings.ReplaceAll(s, "\t", `\t`)
	return s
}
