package sessions

import "context"

type Repository interface {
	Create(ctx context.Context, sessionID string, userID string, expire uint64) error
	Get(ctx context.Context, sessionID string) (string, error) // return UserID
	Delete(ctx context.Context, sessionID string) error
}
