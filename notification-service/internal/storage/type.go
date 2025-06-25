package storage

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// Storage определяет интерфейс хранилища.
type Storage[T any] interface {
	// Store записывает значение в хранилище
	Store(ctx context.Context, key string, value T) error

	// PickAll достает все значения из хранилища
	//
	// convert - функция, конвертирующая значение, содержащееся по ключу в необходимый тип.
	PickAll(ctx context.Context, key string, convert func(string) (T, error)) ([]T, error)
}

// RedisStorage определяет структуру хранилища в Redis.
// Имплементирует интерфейс Storage.
type RedisStorage[T any] struct {
	client *redis.Client
}
