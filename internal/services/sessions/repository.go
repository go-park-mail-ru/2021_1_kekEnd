package sessions

import (
	"time"
)

// Repository go:generate mockgen -destination=mocks/repository_mock.go -package=mocks . Repository
type Repository interface {
	Create(sessionID string, userID string, expire time.Duration) error
	Get(sessionID string) (string, error)
	Delete(sessionID string) error
}
