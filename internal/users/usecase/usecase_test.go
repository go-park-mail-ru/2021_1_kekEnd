package usecase

import (
	"errors"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	mock "github.com/go-park-mail-ru/2021_1_kekEnd/internal/users/repository/mock"
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

	t.Run("CreateUser", func(t *testing.T) {
		repo.On("GetUserByUsername", user.Username).Return(nil, errors.New("user not found")).Once()
		repo.On("CreateUser", user).Return(nil)
		err := uc.CreateUser(user)
		assert.NoError(t, err)
	})

	t.Run("LoginUser", func(t *testing.T) {
		repo.On("GetUserByUsername", user.Username).Return(user, nil)
		repo.On("CheckPassword", user.Password, user).Return(true, nil)
		success := uc.Login(user.Username, user.Password)
		assert.True(t, success)
	})

	t.Run("GetUser", func(t *testing.T) {
		gotUser, err := uc.GetUser(user.Username)
		assert.NoError(t, err)
		assert.Equal(t, user, gotUser)
	})

	t.Run("UpdateUser", func(t *testing.T) {
		updatedUser := models.User{
			Username:      "let_robots_reign",
			Email:         "corrected@ya.ru",
			Password:      "1234567",
		}
		repo.On("UpdateUser", user, updatedUser).Return(&updatedUser, nil)
		_, err := uc.UpdateUser(user, updatedUser)
		assert.NoError(t, err)
	})
}
