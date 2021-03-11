package sessions

import (
	"time"
)

type Repository interface {
	Create(sessionID string, userID string, expire time.Duration) error
	Get(sessionID string) (string, error)
	Delete(sessionID string) error
}
