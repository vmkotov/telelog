package telelog

import (
	"time"
)

func GetCurrentTimestamp() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func TruncateText(text string, maxLength int) string {
	if len(text) <= maxLength {
		return text
	}
	return text[:maxLength] + "..."
}

func SafeString(s string) string {
	if s == "" {
		return "(пусто)"
	}
	return s
}
