package client

import (
	"context"

	"github.com/child6yo/logger-bot/notification-service/internal/storage"
	"github.com/go-telegram/bot"
)

// TelegramBot определяет структуру клиента телеграм бота.
type TelegramBot struct {
	Bot *bot.Bot

	storage storage.Storage[int64]

	ctx    context.Context
	cancel context.CancelFunc
}
