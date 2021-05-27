package sessions

import (
	"time"
)

// UseCase go:generate mockgen -destination=mocks/usecase_mock.go -package=mocks . UseCase
type UseCase interface {
	Create(userID string, expires time.Duration) (string, error)
	GetUser(sessionID string) (string, error)
	Delete(sessionID string) error
}
