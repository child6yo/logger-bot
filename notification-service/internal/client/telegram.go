package client

import (
	"context"

	"github.com/child6yo/logger-bot/notification-service/internal/storage"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// New создает новый экземпляр клиента для телеграм-бота.
func New(token string, storage storage.Storage[int64]) (*TelegramBot, error) {
	b, err := bot.New(token)

	return &TelegramBot{Bot: b, storage: storage}, err
}

// Start запускает работу клиента.
func (tb *TelegramBot) Start() {
	tb.Bot.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, func(ctx context.Context, b *bot.Bot, update *models.Update) {
		chatID := update.Message.Chat.ID
		tb.storage.Store(ctx, storage.ChatIDSet, chatID)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Скоро тут появятся логи...",
		})
	})

	ctx, cancel := context.WithCancel(context.Background())
	tb.ctx = ctx
	tb.cancel = cancel

	tb.Bot.Start(ctx)
}

// Stop останавливает работу клиента.
func (tb *TelegramBot) Stop() {
	tb.cancel()
}
