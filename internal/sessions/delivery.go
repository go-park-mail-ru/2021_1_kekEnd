package sessions

import (
	"time"
)

type Delivery interface {
	Create(userID string, expires time.Duration) (string, error)
	GetUser(sessionID string) (string, error)
	Delete(sessionID string) error
}
