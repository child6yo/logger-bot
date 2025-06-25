package config

import (
	"log"
	"os"
)

// Config - структура, определяющая конфигурацию приложения.
type Config struct {
	KafkaBrokers []string
	KafkaTopics  []string

	BotToken string

	RedisAddress  string
	RedisPassword string
}

// InitConfig инициализирует конфиг (.env).
func InitConfig() Config {
	cfg := Config{}

	cfg.KafkaBrokers = []string{getEnv("KAFKA_BROKER", "localhost:9092")}
	cfg.KafkaTopics = []string{getEnv("KAFKA_TOPIC", "logs")}

	cfg.BotToken = getEnv("BOT_TOKEN", "")

	cfg.RedisAddress = getEnv("REDDIS_ADDRESS", "localhost:6379")
	cfg.RedisPassword = getEnv("REDDIS_PASSWORD", "Qwerty")

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
