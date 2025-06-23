package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/child6yo/logger-bot/internal/app"
	"github.com/child6yo/logger-bot/internal/config"
)

func main() {
	// инициализация конфигурации
	cfg := config.InitConfig()

	// создание экземпляра приложения
	app := app.NewApplication(cfg)

	// запуск приложения
	app.StartApplication()

	// получение сигнала на остановку приложения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	// остановка приложения
	app.StopApplication()
}
