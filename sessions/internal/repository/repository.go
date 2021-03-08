package repository

import (
	"context"
	"github.com/go-redis/redis/v8"
)

type Redis struct {
	client redis.Client
}

func (r *Redis) GetUserSession(ctx context.Context, sessionID string) (string, error) {
	res, err := r.client.Get(ctx, sessionID).Result()
	if err != nil {
		return "", err
	}

	return res, nil
}
