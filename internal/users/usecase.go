package users

import (
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
)

type UseCase interface {
	CreateUser(user *models.User) error

	Login(login, password string) bool

	GetUser(id string) (*models.User, error)

	UpdateUser(id string, newUser *models.User) error
}
