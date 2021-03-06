package auth

import (
	"context"
	"github.com/go-park-mail-ru/2021_1_kekEnd/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error

	GetUserByLoginPassword(ctx context.Context, login, password string) (*models.User, error)

	GetUserByID(ctx context.Context, id int) (*models.User, error)

	UpdateUser(ctx context.Context, id int, newUser *models.User) error
}
