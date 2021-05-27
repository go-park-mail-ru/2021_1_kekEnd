package sessions

import (
	"time"

	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/services/sessions"
	uuid "github.com/satori/go.uuid"
)

// UseCase структура usecase авторизации
type UseCase struct {
	Repository sessions.Repository
}

func addPrefix(id string) string {
	return "sessions:" + id
}

// NewUseCase инициализация структуры хендлера авторизации
func NewUseCase(repo sessions.Repository) *UseCase {
	return &UseCase{
		Repository: repo,
	}
}

// Create создание сессии
func (uc *UseCase) Create(userID string, expires time.Duration) (string, error) {
	sessionID := uuid.NewV4().String()
	sID := addPrefix(sessionID)
	err := uc.Repository.Create(sID, userID, expires)

	return sessionID, err
}

// GetUser получение юзера
func (uc *UseCase) GetUser(sessionID string) (string, error) {
	sID := addPrefix(sessionID)
	return uc.Repository.Get(sID)
}

// Delete удалени сессии
func (uc *UseCase) Delete(sessionID string) error {
	sID := addPrefix(sessionID)
	return uc.Repository.Delete(sID)
}
