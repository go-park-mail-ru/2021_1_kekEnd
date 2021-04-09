package actors

import "github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"

type Repository interface {
	CreateActor(actor models.Actor) error
	EditActor(change models.Actor) (models.Actor, error)
	GetActor(id string) (models.Actor, error)
}
