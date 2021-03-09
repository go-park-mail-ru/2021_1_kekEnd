package sessions

import (
	"context"
	"time"
)

type Delivery interface {
	Create(ctx context.Context, userID string, expires time.Duration) (string, error)
	GetUser(ctx context.Context, sessionID string) (string, error)
	Delete(ctx context.Context, sessionID string) error
}
