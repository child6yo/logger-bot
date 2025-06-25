package storage

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

const (
	ChatIDSet = "chat_ids"
)

// NewRedis создает новое подключение к Redis.
// Пытается подключится к базе данных, до тех пор, пока этого не выйдет.
func NewRedis(address, password string) *redis.Client {
	for {
		rdb := redis.NewClient(&redis.Options{
			Addr:     address,  // адрес Redis
			Password: password, // пароль
			DB:       0,        // номер базы данных
		})

		status := rdb.Ping(context.Background())
		if status.Err() == nil {
			return rdb
		}
	}
}

// NewInt64RedisStorage создает новый экземпляр хранилища для типа int64.
func NewInt64RedisStorage(c *redis.Client) *RedisStorage[int64] {
	return &RedisStorage[int64]{client: c}
}

// Store записывает значение в хранилище
func (rs *RedisStorage[T]) Store(ctx context.Context, key string, value T) error {
	err := rs.client.SAdd(ctx, key, value).Err()

	return fmt.Errorf("storage: failed to store value in %s: %v", key, err)
}

// PickAll достает все значения из хранилища
func (rs *RedisStorage[T]) PickAll(ctx context.Context, key string, convert func(string) (T, error)) (_ []T, err error) {
	values, err := rs.client.SMembers(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("storage: failed to pick values: %v", err)
	}

	result := make([]T, len(values))
	for i, v := range values {
		result[i], err = convert(v)
		if err != nil {
			return nil, fmt.Errorf("storage: failed to convert values: %v", err)
		}
	}

	return result, nil
}
