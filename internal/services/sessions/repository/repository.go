package sessions

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisRepository структура репозитрия авторизации
type RedisRepository struct {
	client *redis.Client
}

// NewRedisRepository инициализация структуры репозитрия авторизации
func NewRedisRepository(rdb *redis.Client) *RedisRepository {
	return &RedisRepository{
		client: rdb,
	}
}

// Create создание сессии
func (r *RedisRepository) Create(sessionID string, userID string, expire time.Duration) error {
	_, err := r.client.Set(context.Background(), sessionID, userID, expire).Result()
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	return nil
}

// Get получение юзера
func (r *RedisRepository) Get(sessionID string) (string, error) {
	userID, err := r.client.Get(context.Background(), sessionID).Result()
	if err != nil {
		return "", fmt.Errorf("session for this user doesn't exits: %w", err)
	}

	return userID, nil
}

// Delete удалени сессии
func (r *RedisRepository) Delete(sessionID string) error {
	_, err := r.client.Del(context.Background(), sessionID).Result()
	if err != nil {
		return fmt.Errorf("failed to delete cookie: %w", err)
	}

	return nil
}
