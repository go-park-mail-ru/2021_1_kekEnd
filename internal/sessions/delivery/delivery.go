package sessions

import (
	"context"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/sessions"
	"time"
)

type Delivery struct {
	UseCase sessions.UseCase
}

func NewDelivery(uc sessions.UseCase) *Delivery {
	return &Delivery{
		UseCase: uc,
	}
}

func (d *Delivery) Create(ctx context.Context, userID string, expires time.Duration) (string, error) {
	sessionsID, err := d.UseCase.Create(ctx, userID, expires)
	if err != nil {
		return "", err
	}

	return sessionsID, nil
}

func (d *Delivery) GetUser(ctx context.Context, sessionsID string) (string, error) {
	user, err := d.UseCase.Check(ctx, sessionsID)
	if err != nil {
		return "", err
	}

	return user, nil
}

func (d *Delivery) Delete(ctx context.Context, sessionsID string) error {
	err := d.UseCase.Delete(ctx, sessionsID)
	if err != nil {
		return err
	}

	return nil
}


