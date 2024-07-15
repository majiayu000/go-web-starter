package utils

import (
	"github.com/majiayu000/gin-starter/internal/models"
	"github.com/sirupsen/logrus"
)

func init() {
	// 配置 logrus
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	logrus.SetLevel(logrus.DebugLevel)
}

func PrintMessages(messages []models.ChatMessage) {
	log := logrus.WithField("context", "chat_messages")

	for i, msg := range messages {
		entry := log.WithFields(logrus.Fields{
			"message_index": i,
			"role":          msg.Role,
		})

		switch msg.Role {
		case "system":
			entry.Info(msg.Content) // 系统消息使用默认颜色（通常是白色）
		case "user":
			entry.WithField("color", "green").Info(msg.Content) // 用户消息使用绿色
		case "assistant":
			entry.WithField("color", "cyan").Info(msg.Content) // 助手消息使用青色
		default:
			entry.Debug(msg.Content) // 其他消息使用默认的调试颜色
		}
	}
}
