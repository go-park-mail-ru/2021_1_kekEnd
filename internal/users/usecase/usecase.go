package usecase

import (
	"errors"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/users"
)

type UsersUseCase struct {
	userRepository users.UserRepository
}

func NewUsersUseCase(repo users.UserRepository) *UsersUseCase {
	return &UsersUseCase{
		userRepository: repo,
	}
}

func (usersUC *UsersUseCase) CreateUser(user *models.User) error {
	_, err := usersUC.userRepository.GetUserByUsername(user.Username)
	if err == nil {
		return errors.New("user already exists")
	}
	return usersUC.userRepository.CreateUser(user)
}

func (usersUC *UsersUseCase) Login(login, password string) bool {
	user, err := usersUC.userRepository.GetUserByUsername(login)
	if err != nil {
		return false
	}
	correct, err := usersUC.userRepository.CheckPassword(password, user)
	if err != nil {
		return false
	}
	return correct
}

func (usersUC *UsersUseCase) GetUser(username string) (*models.User, error) {
	return usersUC.userRepository.GetUserByUsername(username)
}

func (usersUC *UsersUseCase) UpdateUser(user *models.User, change models.User) (*models.User, error) {
	err := usersUC.userRepository.CheckEmailUnique(change.Email)
	if err != nil {
		return nil, err
	}

	return usersUC.userRepository.UpdateUser(user, change)
}
