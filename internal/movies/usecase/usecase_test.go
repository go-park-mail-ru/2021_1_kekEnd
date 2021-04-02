package usecase

import (
	"errors"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	mock "github.com/go-park-mail-ru/2021_1_kekEnd/internal/movies/repository/mock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMoviesUseCase(t *testing.T) {
	repo := &mock.MoviesStorageMock{}
	uc := NewMoviesUseCase(repo)

	movie := &models.Movie{
		ID:          "7",
		Title:       "Yet another mock",
		Description: "Generic description",
	}

	t.Run("CreateMovie", func(t *testing.T) {
		repo.On("GetMovieByID", movie.ID).Return(nil, errors.New("movie not found")).Once()
		repo.On("CreateMovie", movie).Return(nil)
		err := uc.CreateMovie(movie)
		assert.NoError(t, err)
	})

	t.Run("GetMovie", func(t *testing.T) {
		repo.On("GetMovieByID", movie.ID).Return(movie, nil)
		gotMovie, err := uc.GetMovie(movie.ID)
		assert.NoError(t, err)
		assert.Equal(t, movie, gotMovie)
	})
}

func TestMoviesUseCaseErrors(t *testing.T) {
	repo := &mock.MoviesStorageMock{}
	uc := NewMoviesUseCase(repo)

	movie := &models.Movie{
		ID:          "7",
		Title:       "Yet another mock",
		Description: "Generic description",
	}

	t.Run("CreateExistingMovie", func(t *testing.T) {
		repo.On("GetMovieByID", movie.ID).Return(movie, nil)
		repo.On("CreateMovie", movie).Return(nil)
		err := uc.CreateMovie(movie)
		assert.Error(t, err)
		assert.Equal(t, "movie already exists", err.Error())
	})

	t.Run("GetMovieError", func(t *testing.T) {
		wrongMovieID := "42"
		repo.On("GetMovieByID", wrongMovieID).Return(nil, errors.New("movie not found"))
		gotMovie, err := uc.GetMovie(wrongMovieID)
		assert.Nil(t, gotMovie)
		assert.Error(t, err)
		assert.Equal(t, "movie not found", err.Error())
	})
}
