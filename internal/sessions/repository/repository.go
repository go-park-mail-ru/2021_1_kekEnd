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

func NewRedisRepository() *RedisRepository {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	return &RedisRepository{
		client: rdb,
	}
}

func (r *RedisRepository) Create(ctx context.Context, sessionID string, userID string, expire time.Duration) error {
	res := r.client.Set(ctx, sessionID, userID, expire)
	return res.Err()
}

func (r *RedisRepository) Get(ctx context.Context, sessionID string) (string, error) {
	userID, err := r.client.Get(ctx, sessionID).Result()
	if err != nil {
		return "", errors.New("session for this user doesn't exits")
	}

	return userID, nil
}

func (r *RedisRepository) Delete(ctx context.Context, sessionID string) error {
	res := r.client.Del(ctx, sessionID)
	return res.Err()
}