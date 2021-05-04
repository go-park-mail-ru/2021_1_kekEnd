package usecase

import (
	"errors"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/users/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUsersUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mocks.NewMockUserRepository(ctrl)
	uc := NewUsersUseCase(repo)

	user := &models.User{
		Username: "let_robots_reign",
		Email:    "sample@ya.ru",
		Password: "1234",
	}

	t.Run("CreateUser", func(t *testing.T) {
		repo.EXPECT().GetUserByUsername(user.Username).Return(nil, errors.New("user not found"))
		repo.EXPECT().CreateUser(user).Return(nil)
		err := uc.CreateUser(user)
		assert.NoError(t, err)
	})

	t.Run("LoginUser", func(t *testing.T) {
		repo.EXPECT().GetUserByUsername(user.Username).Return(user, nil)
		repo.EXPECT().CheckPassword(user.Password, user).Return(true, nil)
		success := uc.Login(user.Username, user.Password)
		assert.True(t, success)
	})

	t.Run("GetUser", func(t *testing.T) {
		repo.EXPECT().GetUserByUsername(user.Username).Return(user, nil)
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
		repo.EXPECT().CheckEmailUnique("corrected@ya.ru").Return(nil)
		repo.EXPECT().UpdateUser(user, updatedUser).Return(&updatedUser, nil)
		_, err := uc.UpdateUser(user, updatedUser)
		assert.NoError(t, err)
	})

	t.Run("Subscribe", func(t *testing.T) {
		sub := "whaeva"
		user := "let_robots_reign"

		repo.EXPECT().CheckUnsubscribed(sub, user).Return(nil, true)
		repo.EXPECT().Subscribe(sub, user).Return(nil)

		err := uc.Subscribe(sub, user)
		assert.NoError(t, err)
	})

	t.Run("Unsubscribe", func(t *testing.T) {
		sub := "whaeva"
		user := "let_robots_reign"

		repo.EXPECT().CheckUnsubscribed(sub, user).Return(nil, false)
		repo.EXPECT().Unsubscribe(sub, user).Return(nil)

		err := uc.Unsubscribe(sub, user)
		assert.NoError(t, err)
	})

	t.Run("GetSubscribers", func(t *testing.T) {
		user := "let_robots_reign"

		repo.EXPECT().GetSubscribers(0, user).Return(0, []*models.UserNoPassword{}, nil)

		_, subs, err := uc.GetSubscribers(1, user)
		assert.NoError(t, err)
		assert.Equal(t, subs, []*models.UserNoPassword{})
	})

	t.Run("GetSubscriptions", func(t *testing.T) {
		user := "let_robots_reign"

		repo.EXPECT().GetSubscriptions(0, user).Return(0, []*models.UserNoPassword{}, nil)

		_, subs, err := uc.GetSubscriptions(1, user)
		assert.NoError(t, err)
		assert.Equal(t, subs, []*models.UserNoPassword{})
	})

}

func TestUsersUseCaseErrors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mocks.NewMockUserRepository(ctrl)
	uc := NewUsersUseCase(repo)

	user := &models.User{
		Username: "let_robots_reign",
		Email:    "sample@ya.ru",
		Password: "1234",
	}

	t.Run("CreateExistingUser", func(t *testing.T) {
		repo.EXPECT().GetUserByUsername(user.Username).Return(user, nil)
		err := uc.CreateUser(user)
		assert.Error(t, err)
		assert.Equal(t, "user already exists", err.Error())
	})

	t.Run("LoginWrongUsername", func(t *testing.T) {
		//repo.EXPECT().CheckPassword(user.Password, user).Return(true, nil)
		repo.EXPECT().GetUserByUsername("nonexistent_user").Return(nil, errors.New("user not found"))
		success := uc.Login("nonexistent_user", user.Password)
		assert.False(t, success)
	})

	t.Run("LoginWrongPassword", func(t *testing.T) {
		wrongPassword := "123"
		repo.EXPECT().GetUserByUsername(user.Username).Return(user, nil)
		repo.EXPECT().CheckPassword(wrongPassword, user).Return(false, nil)
		success := uc.Login(user.Username, wrongPassword)
		assert.False(t, success)
	})

	t.Run("GetUserError", func(t *testing.T) {
		wrongUsername := "nonexistent_user"
		repo.EXPECT().GetUserByUsername(wrongUsername).Return(nil, errors.New("user not found"))
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

		repo.EXPECT().CheckEmailUnique(update.Email).Return(errors.New("user not found"))
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

		repo.EXPECT().CheckEmailUnique(update.Email).Return(errors.New("user not found"))
		_, err := uc.UpdateUser(user, update)
		assert.Error(t, err)
		assert.Equal(t, "user not found", err.Error())
	})

	t.Run("SubscribeTwice", func(t *testing.T) {
		sub := "whaeva"
		user := "let_robots_reign"

		repo.EXPECT().CheckUnsubscribed(sub, user).Return(nil, false)

		err := uc.Subscribe(sub, user)
		assert.Error(t, err)
	})

	t.Run("UnsubscribeTwice", func(t *testing.T) {
		sub := "whaeva"
		user := "let_robots_reign"

		repo.EXPECT().CheckUnsubscribed(sub, user).Return(nil, true)

		err := uc.Unsubscribe(sub, user)
		assert.Error(t, err)
	})
}
