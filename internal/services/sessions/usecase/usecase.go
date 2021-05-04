package sessions

import (
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/services/sessions"
	"github.com/satori/go.uuid"
	"time"
)

type UseCase struct {
	Repository sessions.Repository
}

func addPrefix(id string) string {
	return "sessions:" + id
}

func NewUseCase(repo sessions.Repository) *UseCase {
	return &UseCase{
		Repository: repo,
	}
}

func (uc *UseCase) Create(userID string, expires time.Duration) (string, error) {
	sessionID := uuid.NewV4().String()
	sID := addPrefix(sessionID)
	err := uc.Repository.Create(sID, userID, expires)

	return sessionID, err
}

func (uc *UseCase) GetUser(sessionID string) (string, error) {
	sID := addPrefix(sessionID)
	return uc.Repository.Get(sID)
}

func (uc *UseCase) Delete(sessionID string) error {
	sID := addPrefix(sessionID)
	return uc.Repository.Delete(sID)
}
