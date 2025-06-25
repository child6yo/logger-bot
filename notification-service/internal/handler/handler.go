package handler

import "github.com/child6yo/logger-bot/notification-service/internal/client"

// NewHandler создает новый экземпляр Handler.
func NewHandler(c *client.TelegramBot) *Handler {
	return &Handler{
		LogsHandler: NewLogsHandler(c),
	}
}
