package actors

import "github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"

//go:generate mockgen -destination=mocks/usecase.go -package=mocks . UseCase
type UseCase interface {
	// CreateActor(user models.User, actor models.Actor) error
	GetActor(id string, username string) (models.Actor, error)
	// EditActor(user models.User, change models.Actor) (models.Actor, error)
	LikeActor(username string, actorID int) error
	UnlikeActor(username string, actorID int) error
}
