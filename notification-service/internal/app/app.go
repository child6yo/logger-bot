package app

import (
	"log"

	"github.com/child6yo/logger-bot/notification-service/internal/client"
	"github.com/child6yo/logger-bot/notification-service/internal/config"
	"github.com/child6yo/logger-bot/notification-service/internal/consumer"
	"github.com/child6yo/logger-bot/notification-service/internal/handler"
	"github.com/child6yo/logger-bot/notification-service/internal/storage"
)

// Application определяет структуру менеджера приложения.
// Управляет его работой.
type Application struct {
	config.Config

	consumer *consumer.Connection
	bot      *client.TelegramBot
}

// NewApplication создает новый экземпляр менеджера приложения.
func NewApplication(cfg config.Config) *Application {
	return &Application{
		Config: cfg,
	}
}

// StartApplication запускает приложение.
func (a *Application) StartApplication() {
	redis := storage.NewRedis(a.RedisAddress, a.RedisPassword)
	storage := storage.NewInt64RedisStorage(redis)

	bot, err := client.New(a.BotToken, storage)
	if err != nil {
		log.Panic(err)
	}
	a.bot = bot

	go func() {
		a.bot.Start()
	}()

	handler := handler.NewHandler(bot, storage)

	a.consumer = consumer.NewConnection(a.KafkaBrokers, a.KafkaTopics, "notification-group", 1, *handler)

	go func() {
		for {
			err := a.consumer.RunConsumers()
			if err == nil {
				break
			}
		}
	}()
}

// StopApplication останавливает работу приложения.
func (a *Application) StopApplication() {
	a.bot.Stop()
	a.consumer.StopConsumers()
}
