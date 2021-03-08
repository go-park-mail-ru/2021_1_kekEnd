package mock

import (
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	"github.com/stretchr/testify/mock"
)

type UserStorageMock struct {
	mock.Mock
}

func (storageMock *UserStorageMock) CreateUser(user *models.User) error {
	args := storageMock.Called(user)
	return args.Error(0)
}

func (storageMock *UserStorageMock) GetUserByUsername(username string) (*models.User, error) {
	args := storageMock.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (storageMock *UserStorageMock) CheckPassword(password string, user *models.User) (bool, error) {
	args := storageMock.Called(password, user)
	return args.Get(0).(bool), args.Error(1)
}

func (storageMock *UserStorageMock) UpdateUser(username string, newUser *models.User) error {
	args := storageMock.Called(username, newUser)
	return args.Error(0)
}
