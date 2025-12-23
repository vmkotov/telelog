package test

import (
	"testing"
	"time"
	
	"github.com/vmkotov/telelog"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func TestSimpleLogger(t *testing.T) {
	logger := telelog.SimpleNew()
	
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
	
	logger.LogMessage(msg, "private")
	logger.LogCommand(msg, "test")
	
	t.Log("✅ Basic test passed")
}

func TestDeployNotification(t *testing.T) {
	logger := telelog.SimpleNew()
	
	deployInfo := map[string]string{
		"version":     "1.0.0",
		"environment": "production",
		"branch":      "main",
		"commit_hash": "abc123",
		"deployer":    "Test Runner",
		"timestamp":   time.Now().Format("2006-01-02 15:04:05"),
	}
	
	logger.SendDeployNotification(deployInfo)
	
	t.Log("✅ Deploy notification test passed")
}
