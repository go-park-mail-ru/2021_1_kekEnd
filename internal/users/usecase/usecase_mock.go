package usecase

import (
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	"github.com/stretchr/testify/mock"
)

type UsersUseCaseMock struct {
	mock.Mock
}

func (mockUC *UsersUseCaseMock) CreateUser(user *models.User) error {
	args := mockUC.Called(user)
	return args.Error(0)
}

func (mockUC *UsersUseCaseMock) Login(login, password string) bool {
	args := mockUC.Called(login, password)
	return args.Get(0).(bool)
}

func (mockUC *UsersUseCaseMock) GetUser(username string) (*models.User, error) {
	args := mockUC.Called(username)
	return args.Get(0).(*models.User), args.Error(1)
}

func (mockUC *UsersUseCaseMock) UpdateUser(username *models.User, newUser models.User) error {
	args := mockUC.Called(username, newUser)
	return args.Error(0)
}
