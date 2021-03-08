package sessions

import (
	"context"
	"github.com/go-park-mail-ru/2021_1_kekEnd/sessions/internal"
	"github.com/satori/go.uuid"
)

type UseCase struct {
	Repository sessions.Repository
}

func addPrefix(id string) string {
	return "sessions:" + id
}

func (uc *UseCase) Create(userID uuid.UUID, expires uint64) (string, error){
	ctx := context.Background()
	sessionID := uuid.NewV4().String()
	sID := addPrefix(sessionID)
	err := uc.Repository.Create(ctx, sID, userID.String(), expires)
	if err != nil {
		return "", err
	}

	return sessionID, nil
}


func (uc *UseCase) Check(sessionID uuid.UUID) (string, error) {
	ctx := context.Background()
	sID := addPrefix(sessionID.String())
	user, err := uc.Repository.Get(ctx, sID)
	if err != nil {
		return "", err
	}

	return user, nil
}

func (us *UseCase) Delete(sessionID uuid.UUID) error {
	ctx := context.Background()
	sID := addPrefix(sessionID.String())
	err := us.Repository.Delete(ctx, sID)
	if err != nil{
		return err
	}

	return nil
}