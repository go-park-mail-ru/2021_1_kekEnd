package sessions

import (
	"context"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/sessions"
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

func (uc *UseCase) Create(ctx context.Context, userID string, expires time.Duration) (string, error){
	sessionID := uuid.NewV4().String()
	sID := addPrefix(sessionID)
	err := uc.Repository.Create(ctx, sID, userID, expires)
	if err != nil {
		return "", err
	}

	return sessionID, nil
}


func (uc *UseCase) Check(ctx context.Context, sessionID string) (string, error) {
	sID := addPrefix(sessionID)
	user, err := uc.Repository.Get(ctx, sID)
	if err != nil {
		return "", err
	}

	return user, nil
}

func (uc *UseCase) Delete(ctx context.Context, sessionID string) error {
	sID := addPrefix(sessionID)
	err := uc.Repository.Delete(ctx, sID)
	if err != nil{
		return err
	}

	return nil
}