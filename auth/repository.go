package auth

import (
	"../models"
	"context"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error

	GetUser(ctx context.Context, login, password string) (*models.User, error)
}
