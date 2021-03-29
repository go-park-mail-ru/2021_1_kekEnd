package localstorage

import (
	"errors"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	"sync"
)

type ActorsLocalStorage struct {
	actors map[string]models.Actor
	mutex  sync.Mutex
}

func NewActorsLocalStorage() *ActorsLocalStorage {
	return &ActorsLocalStorage{
		actors: make(map[string]models.Actor),
	}
}

func (a *ActorsLocalStorage) CreateActor(actor models.Actor) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	a.actors[actor.ID] = actor

	return nil
}

func (a *ActorsLocalStorage) GetActor(id string) (models.Actor, error) {
	actor, exists := a.actors[id]
	if !exists {
		return models.Actor{}, errors.New("actor doesn't exist")
	}

	return actor, nil
}

func (a *ActorsLocalStorage) EditActor(change models.Actor) (models.Actor, error) {
	actor, exists := a.actors[change.ID]
	if !exists {
		return models.Actor{}, errors.New("actor doesn't exist")
	}


	return a.actors[actor.ID], nil
}