package usecase

import (
	"errors"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/movies"
)

type MoviesUseCase struct {
	movieRepository movies.MovieRepository
}

func NewMoviesUseCase(repo movies.MovieRepository) *MoviesUseCase {
	return &MoviesUseCase{
		movieRepository: repo,
	}
}

func (moviesUC *MoviesUseCase) CreateMovie(movie *models.Movie) error {
	_, err := moviesUC.movieRepository.GetMovieByID(movie.ID)
	if err == nil {
		return errors.New("movie already exists")
	}
	return moviesUC.movieRepository.CreateMovie(movie)
}

func (moviesUC *MoviesUseCase) GetMovie(id string) (*models.Movie, error) {
	return moviesUC.movieRepository.GetMovieByID(id)
}

func (moviesUC *MoviesUseCase) UpdateMovie(id string, newMovie *models.Movie) error {
	return moviesUC.movieRepository.UpdateMovie(id, newMovie)
}
