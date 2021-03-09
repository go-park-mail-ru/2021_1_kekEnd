package localstorage

import (
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserLocalStorage(t *testing.T) {
	storage := NewUserLocalStorage()

	user := &models.User{
		Username:      "let_robots_reign",
		Email:         "sample@ya.ru",
		Password:      "1234",
		MoviesWatched: 4,
		ReviewsNumber: 2,
	}

	err := storage.CreateUser(user)
	assert.NoError(t, err)

	gotUser, err := storage.GetUserByUsername("let_robots_reign")
	assert.NoError(t, err)
	assert.Equal(t, user, gotUser)

	gotUser, err = storage.GetUserByUsername("unknown")
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "user not found")

	checkPass, err := storage.CheckPassword("1234", user)
	assert.NoError(t, err)
	assert.True(t, checkPass)

	err = storage.UpdateUser("let_robots_reign", &models.User{
		Username:      "let_robots_reign",
		Email:         "corrected@ya.ru",
		Password:      "12345",
	})
	assert.NoError(t, err)
	updatedUser, err := storage.GetUserByUsername("let_robots_reign")
	assert.NoError(t, err)
	assert.Equal(t, updatedUser.Email, "corrected@ya.ru")
	assert.Equal(t, updatedUser.Password, "12345")
}
