package app

import (
	"log"

	"github.com/child6yo/logger-bot/scrapper-service/internal/config"
	"github.com/child6yo/logger-bot/scrapper-service/internal/producer"
	"github.com/child6yo/logger-bot/scrapper-service/internal/service"
)

// Application определяет структуру менеджера приложения.
// Управляет его работой.
type Application struct {
	config.Config

	producer *producer.KafkaProducer
}

// NewApplication создает новый экземпляр менеджера приложения.
func NewApplication(cfg config.Config) *Application {
	return &Application{
		Config: cfg,
	}
}

// StartApplication запускает приложение.
func (a *Application) StartApplication() {
	a.producer = producer.NewKafkaProducer(a.KafkaBrokers, a.KafkaTopic)
	err := a.producer.StartProducer()
	if err != nil {
		log.Fatal(err)
	}

	scrapper := service.NewScraperService(a.producer)
	go func() {
		err := scrapper.Start(a.LogFilepath, a.LogFilter)
		if err != nil {
			log.Fatal(err)
		}
	}()
}

// StopApplication останавливает работу приложения.
func (a *Application) StopApplication() {
	a.producer.StopProducer()
}
