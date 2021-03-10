package sessions

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisRepository struct {
	client *redis.Client
}

func NewRedisRepository(rdb *redis.Client) *RedisRepository {
	return &RedisRepository{
		client: rdb,
	}
}

func (r *RedisRepository) Create(sessionID string, userID string, expire time.Duration) error {
	_, err := r.client.Set(context.Background(), sessionID, userID, expire).Result()
	return fmt.Errorf("failed to create session: %w", err)
}

func (r *RedisRepository) Get(sessionID string) (string, error) {
	userID, err := r.client.Get(context.Background(), sessionID).Result()
	if err != nil {
		return "", fmt.Errorf("session for this user doesn't exits: %w", err)
	}

	return userID, nil
}

func (r *RedisRepository) Delete(sessionID string) error {
	_, err := r.client.Del(context.Background(), sessionID).Result()
	if err != nil {
		return fmt.Errorf("failed to delete cookie: %w", err)
	}

	return nil
}
