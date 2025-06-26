package handler

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/child6yo/logger-bot/notification-service/internal/client"
	"github.com/child6yo/logger-bot/notification-service/internal/storage"
	"github.com/go-telegram/bot"
)

// NewLogsHandler создает новый экземпляр Logs.
func NewLogsHandler(bot *client.TelegramBot, storage storage.Storage[int64]) *Logs {
	return &Logs{bot: bot, storage: storage}
}

// Handle отправляет логи в чат.
func (l *Logs) Handle(ctx context.Context, message []byte) error {
	convert := func(s string) (int64, error) {
		n, err := strconv.Atoi(s)
		if err != nil {
			return 0, fmt.Errorf("failed to convert %s to int", s)
		}
		return int64(n), nil
	}

	chats, err := l.storage.PickAll(ctx, storage.ChatIDSet, convert)
	if err != nil {
		return fmt.Errorf("handler: failed to pick chats: %v", err)
	}

	for _, chat := range chats {
		_, err := l.bot.Bot.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chat,
			Text:   string(message),
		})

		if err != nil {
			log.Printf("[ERROR] handler: failed to send message to %d: %v", chat, err)
		}
	}

	return nil
}
