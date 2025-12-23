package telelog

import (
	"fmt"
	"log"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (tl *loggerImpl) sendDeployToChat(deployInfo map[string]string) {
	if !tl.enabled || tl.bot == nil {
		log.Println("⚠️ TeleLogger не инициализирован для отправки уведомлений о деплое")
		return
	}

	text := tl.formatDeployMessage(deployInfo)

	msg := tgbotapi.NewMessage(tl.logChatID, text)
	
	
	if _, err := tl.bot.Send(msg); err != nil {
		log.Printf("❌ Не удалось отправить уведомление о деплое: %v", err)
	} else {
		log.Printf("✅ Уведомление о деплое отправлено в чат %d", tl.logChatID)
	}
}

func (tl *loggerImpl) formatDeployMessage(deployInfo map[string]string) string {
	version := deployInfo["version"]
	if version == "" {
		version = "unknown"
	}

	commitHash := deployInfo["commit_hash"]
	if commitHash == "" {
		commitHash = "unknown"
	}

	branch := deployInfo["branch"]
	if branch == "" {
		branch = "unknown"
	}

	deployer := deployInfo["deployer"]
	if deployer == "" {
		deployer = "unknown"
	}

	environment := deployInfo["environment"]
	if environment == "" {
		environment = "production"
	}

	timestamp := deployInfo["timestamp"]
	if timestamp == "" {
		timestamp = time.Now().Format("2006-01-02 15:04:05")
	}

	var additionalInfo strings.Builder
	for key, value := range deployInfo {
		if key == "version" || key == "commit_hash" || key == "branch" ||
		   key == "deployer" || key == "environment" || key == "timestamp" {
			continue
		}
		additionalInfo.WriteString(fmt.Sprintf("• %s: %s\n", key, value))
	}

	additionalText := additionalInfo.String()
	if additionalText != "" {
		additionalText = "\n📊 Дополнительно:\n" + additionalText
	}

	botUsername := "unknown"
	if tl.botUsername != "" {
		botUsername = tl.botUsername
	}

	return fmt.Sprintf(
		"🚀 УВЕДОМЛЕНИЕ О ДЕПЛОЕ\n\n"+
			"📦 Версия: %s\n"+
			"🔧 Окружение: %s\n"+
			"🌿 Ветка: %s\n"+
			"📝 Коммит: %s\n"+
			"👤 Деплойер: %s\n"+
			"⏰ Время: %s\n"+
			"🤖 Бот: @%s\n"+
			"%s\n"+
			"✅ Деплой успешно завершен!",
		version,
		environment,
		branch,
		commitHash,
		deployer,
		timestamp,
		botUsername,
		additionalText,
	)
}
