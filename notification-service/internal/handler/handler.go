package handler

import (
	"github.com/child6yo/logger-bot/notification-service/internal/client"
	"github.com/child6yo/logger-bot/notification-service/internal/storage"
)

// NewHandler создает новый экземпляр Handler.
func NewHandler(c *client.TelegramBot, storage storage.Storage[int64]) *Handler {
	return &Handler{
		LogsHandler: NewLogsHandler(c, storage),
	}
}
