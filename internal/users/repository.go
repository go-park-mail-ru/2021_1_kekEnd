package users

import (
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
)

type UserRepository interface {
	CreateUser(user *models.User) error

	GetUserByLoginPassword(login, password string) (*models.User, error)

	GetUserByID(id int) (*models.User, error)

	UpdateUser(id int, newUser *models.User) error
}
