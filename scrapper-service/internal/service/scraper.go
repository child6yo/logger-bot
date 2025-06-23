package service

import (
	"fmt"
	"io"
	"log"
	"regexp"

	"github.com/child6yo/logger-bot/internal/producer"
	"github.com/nxadm/tail"
)

// ScraperService имплементирует интерфейс Scraper.
type ScraperService struct {
	producer producer.Producer
}

// NewScraperService создает новый экземпляр ScraperService.
func NewScraperService() *ScraperService {
	return &ScraperService{}
}

// Start запускает работу сборщика логов.
// На вход принимает два параметра:
//   - filepath - путь к файлу, с которого будут читаться логи.
//   - filter - регулярное выражение, фильтрующее необходимые логи.
func (ss *ScraperService) Start(filepath, filter string) error {
	
	// nxadm/tail следит за файлов в live-формате, чтение начинается с последней строки и ждет новых
	t, err := tail.TailFile(filepath, tail.Config{
		Follow:   true,
		Location: &tail.SeekInfo{Offset: 0, Whence: io.SeekEnd},
	})

	if err != nil {
		return fmt.Errorf("failed to start scrapper: %v", err)
	}

	for line := range t.Lines {
		if ok, _ := regexp.MatchString(filter, line.Text); ok {
			err := ss.producer.SendStringMessage(line.Text)
			if err != nil {
				log.Printf("[ERROR] scraper service: %v", err)
			}
		}
	}

	return nil
}
