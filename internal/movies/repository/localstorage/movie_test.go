package localstorage

import (
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMovieLocalStorage(t *testing.T) {
	storage := NewMovieLocalStorage()

	movie := &models.Movie{
		ID:          "7",
		Title:       "Yet another movie",
		Description: "Generic description",
	}

	t.Run("CreateMovie", func(t *testing.T) {
		err := storage.CreateMovie(movie)
		assert.NoError(t, err)
	})

	t.Run("SuccessfulGetMovie", func(t *testing.T) {
		gotMovie, err := storage.GetMovieByID("7")
		assert.NoError(t, err)
		assert.Equal(t, movie, gotMovie)
	})

	t.Run("UnsuccessfulGetMovie", func(t *testing.T) {
		gotMovie, err := storage.GetMovieByID("42")
		assert.Nil(t, gotMovie)
		assert.Error(t, err)
		assert.Equal(t, "movie not found", err.Error())
	})
}
