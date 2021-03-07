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
	_, err := usersUC.userRepository.GetUserByLoginPassword(login, password)
	if err != nil {
		return false
	}
	return true
}

func (usersUC *UsersUseCase) GetUser(id int) (*models.User, error) {
	user, err := usersUC.userRepository.GetUserByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (usersUC *UsersUseCase) UpdateUser(id int, newUser *models.User) error {
	if err := usersUC.userRepository.UpdateUser(id, newUser); err != nil {
		return err
	}
	return nil
}
