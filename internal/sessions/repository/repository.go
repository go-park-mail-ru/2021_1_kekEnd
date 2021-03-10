package sessions

import (
	"context"
	"errors"
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
	ctx := context.Background()
		_, err := r.client.Set(ctx, sessionID, userID, expire).Result()
	return err
}

func (r *RedisRepository) Get(sessionID string) (string, error) {
	ctx := context.Background()
	userID, err := r.client.Get(ctx, sessionID).Result()
	if err != nil {
		return "", errors.New("session for this user doesn't exits")
	}

	return userID, nil
}

func (r *RedisRepository) Delete(sessionID string) error {
	ctx := context.Background()
	_, err := r.client.Del(ctx, sessionID).Result()
	if err != nil {
		return errors.New("failed to delete cookie")
	}

	return nil
}