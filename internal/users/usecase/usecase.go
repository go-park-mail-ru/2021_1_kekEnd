package usecase

import (
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
	return usersUC.userRepository.CreateUser(user)
}

func (usersUC *UsersUseCase) Login(login, password string) bool {
	user, err := usersUC.userRepository.GetUserByUsername(login)
	if err != nil {
		return false
	}
	if password == user.Password {
		return true
	}
	return false
}

func (usersUC *UsersUseCase) GetUser(username string) (*models.User, error) {
	return usersUC.userRepository.GetUserByUsername(username)
}

func (usersUC *UsersUseCase) UpdateUser(id string, newUser *models.User) error {
	if err := usersUC.userRepository.UpdateUser(id, newUser); err != nil {
		return err
	}
	return nil
}
