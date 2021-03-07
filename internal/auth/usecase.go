package auth

import (
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
)

type UseCase interface {
	CreateUser(user *models.User) error

	Login(login, password string) bool

	GetUser(id int) (*models.User, error)

	UpdateUser(id int, newUser *models.User) error
}
