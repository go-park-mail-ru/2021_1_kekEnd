package usecase

import (
	"context"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/auth"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
)

type AuthUseCase struct {
	userRepository auth.UserRepository
}

func NewAuthUseCase(repo auth.UserRepository) *AuthUseCase {
	return &AuthUseCase{
		userRepository: repo,
	}
}

func (authUC *AuthUseCase) SignUp(ctx context.Context, username, firstName, lastName, email, password string) error {
	// хэширование и соль

	user := &models.User{
		Username:      username,
		Email:         email,
		Password:      password,
		FirstName:     firstName,
		LastName:      lastName,
		MoviesWatched: 0,
		ReviewsNumber: 0,
	}

	return authUC.userRepository.CreateUser(ctx, user)
}

func (authUC *AuthUseCase) Login(ctx context.Context, login, password string) bool {
	_, err := authUC.userRepository.GetUserByLoginPassword(ctx, login, password)
	if err != nil {
		return false
	}
	return true
}

func (authUC *AuthUseCase) GetUser(ctx context.Context, id int) (*models.User, error) {
	user, err := authUC.userRepository.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (authUC *AuthUseCase) UpdateUser(ctx context.Context, id int, newUser *models.User) error {
	if err := authUC.userRepository.UpdateUser(ctx, id, newUser); err != nil {
		return err
	}
	return nil
}
