package sessions

import (
	"github.com/go-park-mail-ru/2021_1_kekEnd/sessions"
	uuid "github.com/satori/go.uuid"
)

type Delivery struct {
	UseCase main.sessions
}

func (d *Delivery) Create(userID uuid.UUID, expires uint64) (string, error) {
	sessionsID, err := d.UseCase.Create(userID, expires)
	if err != nil {
		return "", err
	}

	return sessionsID, nil
}

func (d *Delivery) Delete(sessionsID uuid.UUID) error {
	err := d.UseCase.Delete(sessionsID)
	if err != nil {
		return err
	}

	return nil
}

func (d *Delivery) Check(sessionsID string) (string, error) {
	sID, _ :=  uuid.FromString(sessionsID)
	user, err := d.UseCase.Check(sID)
	if err != nil {
		return "", err
	}

	return user, nil
}