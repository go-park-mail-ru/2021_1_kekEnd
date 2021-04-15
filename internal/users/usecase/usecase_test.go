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
		Username: "let_robots_reign",
		Email:    "sample@ya.ru",
		Password: "1234",
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
			Username: "let_robots_reign",
			Email:    "corrected@ya.ru",
			Password: "1234567",
		}
		repo.On("CheckEmailUnique", "corrected@ya.ru").Return(nil)
		repo.On("UpdateUser", user, updatedUser).Return(&updatedUser, nil)
		_, err := uc.UpdateUser(user, updatedUser)
		assert.NoError(t, err)
	})
}

func TestUsersUseCaseErrors(t *testing.T) {
	repo := &mock.UserStorageMock{}
	uc := NewUsersUseCase(repo)

	user := &models.User{
		Username: "let_robots_reign",
		Email:    "sample@ya.ru",
		Password: "1234",
	}

	t.Run("CreateExistingUser", func(t *testing.T) {
		repo.On("GetUserByUsername", user.Username).Return(user, nil)
		repo.On("CreateUser", user).Return(nil)
		err := uc.CreateUser(user)
		assert.Error(t, err)
		assert.Equal(t, "user already exists", err.Error())
	})

	t.Run("LoginWrongUsername", func(t *testing.T) {
		repo.On("CheckPassword", user.Password, user).Return(true, nil).Once()
		repo.On("GetUserByUsername", "nonexistent_user").Return(nil, errors.New("user not found"))
		success := uc.Login("nonexistent_user", user.Password)
		assert.False(t, success)
	})

	t.Run("LoginWrongPassword", func(t *testing.T) {
		wrongPassword := "123"
		repo.On("CheckPassword", wrongPassword, user).Return(false, nil).Once()
		success := uc.Login(user.Username, wrongPassword)
		assert.False(t, success)
	})

	t.Run("GetUserError", func(t *testing.T) {
		wrongUsername := "nonexistent_user"
		repo.On("GetUserByUsername", wrongUsername).Return(nil, errors.New("user not found"))
		gotUser, err := uc.GetUser(wrongUsername)
		assert.Nil(t, gotUser)
		assert.Error(t, err)
		assert.Equal(t, "user not found", err.Error())
	})

	t.Run("UpdateNonExistentUser", func(t *testing.T) {
		nonExistentUser := &models.User{
			Username: "nonexistent_user",
			Email:    "corrected@ya.ru",
			Password: "12345",
		}

		update := models.User{
			Username: "nonexistent_user",
			Email:    "new_email@ya.ru",
			Password: "qwerty",
		}

		repo.On("CheckEmailUnique", "new_email@ya.ru").Return(errors.New("user not found"))
		repo.On("UpdateUser", nonExistentUser, update).Return(nil, errors.New("user not found"))
		_, err := uc.UpdateUser(nonExistentUser, update)
		assert.Error(t, err)
		assert.Equal(t, "user not found", err.Error())
	})

	t.Run("UpdateWrongUser", func(t *testing.T) {
		update := models.User{
			Username: "nonexistent_user",
			Email:    "new_email@ya.ru",
			Password: "qwerty",
		}

		repo.On("CheckEmailUnique", "new_email@ya.ru").Return(errors.New("user not found"))
		repo.On("UpdateUser", user, update).Return(nil, errors.New("user not found"))
		_, err := uc.UpdateUser(user, update)
		assert.Error(t, err)
		assert.Equal(t, "user not found", err.Error())
	})
}
