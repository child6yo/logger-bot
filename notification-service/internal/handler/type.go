package handler

import (
	"context"

	"github.com/child6yo/logger-bot/notification-service/internal/client"
	"github.com/child6yo/logger-bot/notification-service/internal/storage"
)

// LogsHandler определяет интерфейс, отвечающий за обработку
// события получения логов через Consumer.
type LogsHandler interface {
	// Handle отправляет логи в чат.
	Handle(ctx context.Context, message []byte) error
}

// Logs имплементирует интерфейс LogsHandler.
type Logs struct {
	bot     *client.TelegramBot
	storage storage.Storage[int64]
}

// Handler определяет общую структуру обработчика событий.
type Handler struct {
	LogsHandler
}
