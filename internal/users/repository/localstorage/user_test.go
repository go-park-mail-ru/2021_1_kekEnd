package localstorage

import (
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserLocalStorage(t *testing.T) {
	storage := NewUserLocalStorage()

	user := &models.User{
		Username:      "let-robots-reign",
		Email:         "sample@ya.ru",
		Password:      "1234",
		MoviesWatched: 4,
		ReviewsNumber: 2,
	}

	err := storage.CreateUser(user)
	assert.NoError(t, err)

	gotUser, err := storage.GetUserByUsername("let-robots-reign")
	assert.NoError(t, err)
	assert.Equal(t, user, gotUser)

	gotUser, err = storage.GetUserByUsername("unknown")
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "user not found")
}
