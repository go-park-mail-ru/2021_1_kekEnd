package usecase

import (
	"errors"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/actors"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	_const "github.com/go-park-mail-ru/2021_1_kekEnd/pkg/const"
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
	permission := false
	for _, admin := range _const.AdminUsers {
		if user.Username == admin {
			permission = true
			break
		}
	}
	if !permission {
		return errors.New("user does not have permission for this action")
	}

	return u.repository.CreateActor(actor)
}

func (u ActorUseCase) GetActor(id string) (models.Actor, error) {
	return u.repository.GetActor(id)
}

func (u ActorUseCase) EditActor(user models.User, change models.Actor) (models.Actor, error) {
	permission := false
	for _, admin := range _const.AdminUsers {
		if user.Username == admin {
			permission = true
		}
	}
	if !permission {
		return models.Actor{}, errors.New("user does not have permission for this action")
	}

	return u.repository.EditActor(change)
}
