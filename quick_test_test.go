package telelog

import (
	"fmt"
	"testing"
	"time"
	
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func TestSimpleLogger(t *testing.T) {
	logger := SimpleNew()
	
	msg := &tgbotapi.Message{
		MessageID: 1,
		From: &tgbotapi.User{
			ID:        123,
			FirstName: "Test",
			UserName:  "test_user",
		},
		Chat: &tgbotapi.Chat{
			ID:    456,
			Title: "Test Chat",
		},
		Text: "Test message",
		Date: int(time.Now().Unix()),
	}
	
	// Эти методы не должны паниковать
	logger.LogMessage(msg, "private")
	logger.LogCommand(msg, "test")
	logger.LogError(msg, fmt.Errorf("test error"), "test context")
	
	deployInfo := map[string]string{
		"version": "1.0.0",
	}
	logger.SendDeployNotification(deployInfo)
	
	t.Log("✅ Quick test passed")
}

func TestInterface(t *testing.T) {
	var _ TeleLogger = (*loggerImpl)(nil)
	t.Log("✅ Interface implementation verified")
}
