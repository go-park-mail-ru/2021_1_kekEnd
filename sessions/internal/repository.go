package sessions

import "context"

type Repository interface {
	Create(sessionID string, userID string) error
	GetUserSession(ctx context.Context, sessionID string) (string, error) // return UserID
	Delete(ctx context.Context, sessionID string) error
}
