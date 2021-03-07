package usecase

import (
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

func (authUC *AuthUseCase) CreateUser(user *models.User) error {
	return authUC.userRepository.CreateUser(user)
}

func (authUC *AuthUseCase) Login(login, password string) bool {
	_, err := authUC.userRepository.GetUserByLoginPassword(login, password)
	if err != nil {
		return false
	}
	return true
}

func (authUC *AuthUseCase) GetUser(id int) (*models.User, error) {
	user, err := authUC.userRepository.GetUserByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (authUC *AuthUseCase) UpdateUser(id int, newUser *models.User) error {
	if err := authUC.userRepository.UpdateUser(id, newUser); err != nil {
		return err
	}
	return nil
}
