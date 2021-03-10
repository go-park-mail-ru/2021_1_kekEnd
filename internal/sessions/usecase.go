package sessions

import (
	"time"
)

type UseCase interface {
	Create(userID string, expires time.Duration) (string, error)
	Check(sessionID string) (string, error)
	Delete(sessionID string) error
}
