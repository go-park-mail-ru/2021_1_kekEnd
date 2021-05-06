package usecase

import (
	"errors"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/movies/mocks"
	userMocks "github.com/go-park-mail-ru/2021_1_kekEnd/internal/users/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMoviesUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockMovieRepository(ctrl)
	usersRepo := userMocks.NewMockUserRepository(ctrl)
	uc := NewMoviesUseCase(repo, usersRepo)

	movie := &models.Movie{
		ID:          "7",
		Title:       "Yet another mock",
		Description: "Generic description",
		Rating:      8.5,
		RatingCount: 1,
		Genre:       []string{"драма"},
	}

	t.Run("CreateMovie", func(t *testing.T) {
		repo.EXPECT().GetMovieByID(movie.ID, "").Return(nil, errors.New("movie not found"))
		repo.EXPECT().CreateMovie(movie).Return(nil)
		err := uc.CreateMovie(movie)
		assert.NoError(t, err)
	})

	t.Run("GetMovie", func(t *testing.T) {
		repo.EXPECT().GetMovieByID(movie.ID, "").Return(movie, nil)
		gotMovie, err := uc.GetMovie(movie.ID, "")
		assert.NoError(t, err)
		assert.Equal(t, movie, gotMovie)
	})

	t.Run("GetBestMovies", func(t *testing.T) {
		repo.EXPECT().GetBestMovies(0, "").Return(1, []*models.Movie{
			movie,
		}, nil)
		const page = 1
		pages, best, err := uc.GetBestMovies(page, "")
		assert.NoError(t, err)
		assert.Equal(t, 1, pages)
		assert.Equal(t, []*models.Movie{movie}, best)
	})

	t.Run("GetAllGenres", func(t *testing.T) {
		repo.EXPECT().GetAllGenres().Return([]string{"драма"}, nil)
		genres, err := uc.GetAllGenres()
		assert.NoError(t, err)
		assert.Equal(t, []string{"драма"}, genres)
	})

	t.Run("GetMoviesByGenres", func(t *testing.T) {
		repo.EXPECT().GetMoviesByGenres([]string{"драма"}, 0, "").Return(1, []*models.Movie{movie}, nil)
		pages, movies, err := uc.GetMoviesByGenres([]string{"драма"}, 1, "")
		assert.NoError(t, err)
		assert.Equal(t, 1, pages)
		assert.Equal(t, []*models.Movie{movie}, movies)
	})
}

func TestMoviesUseCaseErrors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockMovieRepository(ctrl)
	usersRepo := userMocks.NewMockUserRepository(ctrl)
	uc := NewMoviesUseCase(repo, usersRepo)

	movie := &models.Movie{
		ID:          "7",
		Title:       "Yet another mock",
		Description: "Generic description",
	}

	t.Run("CreateExistingMovie", func(t *testing.T) {
		repo.EXPECT().GetMovieByID(movie.ID, "").Return(movie, nil)
		err := uc.CreateMovie(movie)
		assert.Error(t, err)
		assert.Equal(t, "movie already exists", err.Error())
	})

	t.Run("GetMovieError", func(t *testing.T) {
		wrongMovieID := "42"
		repo.EXPECT().GetMovieByID(wrongMovieID, "").Return(nil, errors.New("movie not found"))
		gotMovie, err := uc.GetMovie(wrongMovieID, "")
		assert.Nil(t, gotMovie)
		assert.Error(t, err)
		assert.Equal(t, "movie not found", err.Error())
	})
}
