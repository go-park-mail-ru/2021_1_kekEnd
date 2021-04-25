package users

import (
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
)

//go:generate mockgen -destination=mocks/usecase.go -package=mocks . UseCase
type UseCase interface {
	CreateUser(user *models.User) error

	Login(login, password string) bool

	GetUser(username string) (*models.User, error)

	UpdateUser(user *models.User, change models.User) (*models.User, error)

	Subscribe(subscriber string, user string) error

	Unsubscribe(subscriber *models.User, user *models.User) error

	GetSubscribers(user *models.User) []models.User

	GetSubscriptions(user *models.User) []models.User
}
