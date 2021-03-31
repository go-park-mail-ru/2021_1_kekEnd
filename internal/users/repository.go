package users

import (
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
)

type UserRepository interface {
	CreateUser(user *models.User) error

	GetUserByUsername(username string) (*models.User, error)

	CheckPassword(password string, user *models.User) (bool, error)

	UpdateUser(user *models.User, change models.User) (*models.User, error)

	CheckEmailUnique(newEmail string) error
}
