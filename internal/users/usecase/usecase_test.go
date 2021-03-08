package usecase

import (
	"errors"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/users/repository/mock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUsersUseCase(t *testing.T) {
	repo := &mock.UserStorageMock{}
	uc := NewUsersUseCase(repo)

	user := &models.User{
		Username:      "let_robots_reign",
		Email:         "sample@ya.ru",
		Password:      "1234",
	}

	repo.On("GetUserByUsername", user.Username).Return(nil, errors.New("user not found")).Once()
	repo.On("CreateUser", user).Return(nil)
	err := uc.CreateUser(user)
	assert.NoError(t, err)

	repo.On("GetUserByUsername", user.Username).Return(user, nil)
	repo.On("CheckPassword", user.Password, user).Return(true, nil)
	success := uc.Login(user.Username, user.Password)
	assert.True(t, success)

	gotUser, err := uc.GetUser(user.Username)
	assert.NoError(t, err)
	assert.Equal(t, user, gotUser)

	updatedUser := &models.User{
		Username:      "let_robots_reign",
		Email:         "corrected@ya.ru",
		Password:      "1234567",
	}
	repo.On("UpdateUser", user.Username, updatedUser).Return(nil)
	err = uc.UpdateUser(user.Username, updatedUser)
	assert.NoError(t, err)
}
