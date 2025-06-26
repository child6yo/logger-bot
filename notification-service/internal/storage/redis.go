package storage

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

const (
	// ChatIDSet - название хранящегося в redis множества айди чатов
	ChatIDSet = "chat_ids"
)

// NewRedis создает новое подключение к Redis.
// Пытается подключится к базе данных, до тех пор, пока этого не выйдет.
func NewRedis(address, password string) *redis.Client {
	for {
		rdb := redis.NewClient(&redis.Options{
			Addr:     address,  // адрес
			Password: password, // пароль
			DB:       0,        // номер базы данных
		})

		status := rdb.Ping(context.Background())
		if status.Err() == nil {
			return rdb
		}

		log.Printf("redis: failed to connect to %s: %v, reconnecting...", status.Err(), address)
	}
}

// NewInt64RedisStorage создает новый экземпляр хранилища для типа int64.
func NewInt64RedisStorage(c *redis.Client) *RedisStorage[int64] {
	return &RedisStorage[int64]{client: c}
}

// Store записывает значение в хранилище
func (rs *RedisStorage[T]) Store(ctx context.Context, key string, value T) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("storage: failed to store value in %s: %v", key, err)
		}
	}()

	err = rs.client.SAdd(ctx, key, value).Err()

	return nil
}

// PickAll достает все значения из хранилища
func (rs *RedisStorage[T]) PickAll(ctx context.Context, key string, convert func(string) (T, error)) ([]T, error) {
	values, err := rs.client.SMembers(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("storage: failed to pick values: %v", err)
	}

	result := make([]T, 0, len(values))
	for _, v := range values {
		val, err := convert(v)
		if err != nil {
			return nil, fmt.Errorf("storage: failed to convert values: %v", err)
		}
		result = append(result, val)
	}

	return result, nil
}
