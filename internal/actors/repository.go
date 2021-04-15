package actors

import "github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"

type Repository interface {
	CreateActor(models.Actor) error
	GetActorByID(id string) (models.Actor, error)
	EditActor(models.Actor) (models.Actor, error)
}
