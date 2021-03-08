package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type SessionsRedis struct {
	client *redis.Client
}

func NewSessionsRedis() *SessionsRedis {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	return &SessionsRedis{
		client: rdb,
	}
}

func (r *SessionsRedis) Create(ctx context.Context, sessionID string, userID string, expire uint64) error {
	res := r.client.Set(ctx, sessionID, userID, 240 * time.Hour)
	return res.Err()
}

func (r *SessionsRedis) Get(ctx context.Context, sessionID string) (string, error) {
	userID, err := r.client.Get(ctx, sessionID).Result()
	if err != nil {
		return "", err
	}

	return userID, nil
}

func (r *SessionsRedis) Delete(ctx context.Context, sessionID string) error {
	res := r.client.Del(ctx, sessionID)
	return res.Err()
}