package sessions

import (
	"context"
	"time"
)

type UseCase interface {
	Create(ctx context.Context, userID string, expires time.Duration) (string, error)
	Check(ctx context.Context, sessionID string) (string, error)
	Delete(ctx context.Context, sessionID string) error
}