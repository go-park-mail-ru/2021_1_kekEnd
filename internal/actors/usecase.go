package actors

import "github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"

type UseCase interface {
	CreateActor(user models.User, actor models.Actor) error
	GetActor(id string) (models.Actor, error)
	EditActor(user models.User, change models.Actor) (models.Actor, error)
}
