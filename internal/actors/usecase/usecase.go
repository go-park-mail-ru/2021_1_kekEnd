package usecase

import (
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/actors"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
)

type ActorUseCase struct {
	repository actors.Repository
}

func NewActorsUseCase(repository actors.Repository) *ActorUseCase {
	return &ActorUseCase{
		repository: repository,
	}
}

func (u ActorUseCase) CreateActor(actor models.Actor) error {
	return u.repository.CreateActor(actor)
}

func (u ActorUseCase) GetActor(id string) (models.Actor, error) {
	return u.repository.GetActor(id)
}

func (u ActorUseCase) EditActor(change models.Actor) (models.Actor, error) {
	return u.repository.EditActor(change)
}
