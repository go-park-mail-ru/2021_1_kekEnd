package usecase

import (
	"context"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/sessions/repository/redis"
	"github.com/satori/go.uuid"
)

type SessionsUseCase struct {
	Repository redis.SessionsRedis
}

func addPrefix(id string) string {
	return "sessions:" + id
}

func NewSessionsUseCase(repo redis.SessionsRedis) *SessionsUseCase {
	return &SessionsUseCase{
		Repository: repo,
	}
}

func (uc *SessionsUseCase) Create(userID uuid.UUID, expires uint64) (string, error){
	ctx := context.Background()
	sessionID := uuid.NewV4().String()
	sID := addPrefix(sessionID)
	err := uc.Repository.Create(ctx, sID, userID.String(), expires)
	if err != nil {
		return "", err
	}

	return sessionID, nil
}


func (uc *SessionsUseCase) Check(sessionID uuid.UUID) (string, error) {
	ctx := context.Background()
	sID := addPrefix(sessionID.String())
	user, err := uc.Repository.Get(ctx, sID)
	if err != nil {
		return "", err
	}

	return user, nil
}

func (uc *SessionsUseCase) Delete(sessionID uuid.UUID) error {
	ctx := context.Background()
	sID := addPrefix(sessionID.String())
	err := uc.Repository.Delete(ctx, sID)
	if err != nil{
		return err
	}

	return nil
}