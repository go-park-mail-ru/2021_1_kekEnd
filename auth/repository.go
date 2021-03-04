package auth

import (
	"context"
	"github.com/go-park-mail-ru/2021_1_kekEnd/kekEnd_main/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error

	GetUser(ctx context.Context, login, password string) (*models.User, error)
}
