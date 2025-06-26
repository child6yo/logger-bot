package service

import "github.com/child6yo/logger-bot/scrapper-service/internal/producer"

// Scraper определяет интерфейс сервиса сборщика логов.
// Сервис собирает логи из файла и направляет их в producer.
type Scraper interface {
	// Start запускает работу сборщика логов.
	// На вход принимает два параметра:
	// 	- filepath - путь к файлу, с которого будут читаться логи.
	//  - filter - регулярное выражение, фильтрующее необходимые логи.
	Start(filepath string, filters []string)
}

// ScraperService имплементирует интерфейс Scraper.
type ScraperService struct {
	producer producer.Producer
}
