package usecase

import (
	"errors"
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

func (u ActorUseCase) CreateActor(user models.User, actor models.Actor) error {
	if !(user.Username == "let_robots_reign" ||
		user.Username == "IfuryI" ||
		user.Username == "grillow" ||
		user.Username == "polyanimal") {
		return errors.New("user does not have permission for this action")
	}

	return u.repository.CreateActor(actor)
}

func (u ActorUseCase) GetActor(id string) (models.Actor, error) {
	return u.repository.GetActor(id)
}

func (u ActorUseCase) EditActor(user models.User, change models.Actor) (models.Actor, error) {
	if !(user.Username == "let_robots_reign" ||
		user.Username == "IfuryI" ||
		user.Username == "grillow" ||
		user.Username == "polyanimal") {
		return models.Actor{}, errors.New("user does not have permission for this action")
	}

	return u.repository.EditActor(change)
}
