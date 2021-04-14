package sessions

import (
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/logger"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/sessions"
	"time"
)

type Delivery struct {
	UseCase sessions.UseCase
	Log *logger.Logger
}

func NewDelivery(uc sessions.UseCase, Log *logger.Logger) *Delivery {
	return &Delivery{
		UseCase: uc,
	}
}

func (d *Delivery) Create(userID string, expires time.Duration) (string, error) {
	return d.UseCase.Create(userID, expires)
}

func (d *Delivery) GetUser(sessionID string) (string, error) {
	return d.UseCase.Check(sessionID)
}

func (d *Delivery) Delete(sessionID string) error {
	return d.UseCase.Delete(sessionID)
}
