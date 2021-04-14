package actors

import "github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"

type UseCase interface {
	CreateActor(actor models.Actor) error
	GetActor(id string) (models.Actor, error)
	EditActor(change models.Actor) error
}
