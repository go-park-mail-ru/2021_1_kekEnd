package auth

import "context"

type UseCase interface {
	SignUp(ctx context.Context, username, email, password string) error

	Login(ctx context.Context, login, password string) (bool, error)
}
