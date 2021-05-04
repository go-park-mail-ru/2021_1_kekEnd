package sessions

import (
	"time"
)

//go:generate mockgen -destination=mocks/delivery_mock.go -package=mocks . Delivery
type Delivery interface {
	Create(userID string, expires time.Duration) (string, error)
	GetUser(sessionID string) (string, error)
	Delete(sessionID string) error
}
