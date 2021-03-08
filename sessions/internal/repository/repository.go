package sessions

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type Redis struct {
	client redis.Client
}

func (r *Redis) Create(ctx context.Context, sessionID string, userID string, expire uint64) error {
	res := r.client.Set(ctx, sessionID, userID, 240 * time.Hour)
	return res.Err()
}

func (r *Redis) Get(ctx context.Context, sessionID string) (string, error) {
	userID, err := r.client.Get(ctx, sessionID).Result()
	if err != nil {
		return "", err
	}

	return userID, nil
}

func (r *Redis) Delete(ctx context.Context, sessionID string) error {
	res := r.client.Del(ctx, sessionID)
	return res.Err()
}