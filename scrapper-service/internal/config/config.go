package config

import (
	"log"
	"os"
)

// Config - структура, определяющая конфигурацию приложения.
type Config struct {
	KafkaBrokers []string
	KafkaTopic   string

	LogFilepath string
	LogFilter   string
}

// InitConfig инициализирует конфиг (.env).
func InitConfig() Config {
	cfg := Config{}

	cfg.KafkaBrokers = []string{getEnv("KAFKA_BROKER", "localhost:9092")}
	cfg.KafkaTopic = getEnv("KAFKA_TOPIC", "logs")

	cfg.LogFilepath = getEnv("LOG_FILEPATH", "logs.log")
	cfg.LogFilter = getEnv("LOG_FILTER", `^.*\[ERROR\].*$`)

	return cfg
}

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		if defaultValue != "" {
			log.Printf("config: failed to load env key = %s, defaul value = %s", key, defaultValue)
		}
		return defaultValue
	}
	return value
}
