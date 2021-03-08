package sessions

import (
	"context"
	"time"
)

type Repository interface {
	Create(ctx context.Context, sessionID string, userID string, expire time.Duration) error
	Get(ctx context.Context, sessionID string) (string, error) // return UserID
	Delete(ctx context.Context, sessionID string) error
}
