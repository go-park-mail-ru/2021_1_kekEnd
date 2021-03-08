package sessions

import uuid "github.com/satori/go.uuid"

type UseCase interface {
	Create(userID uuid.UUID, expires uint64) (string, error)
	Check(sessionID uuid.UUID) (string, error)
	Delete(sessionID uuid.UUID) error
}