package usecase

import (
	"errors"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/movies"
	_const "github.com/go-park-mail-ru/2021_1_kekEnd/pkg/const"
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

func (moviesUC *MoviesUseCase) GetBestMovies(page int) (int, []*models.Movie, error) {
	startIndex := (page - 1) * _const.MoviesPageSize
	return moviesUC.movieRepository.GetBestMovies(startIndex)
}

func (moviesUC *MoviesUseCase) GetMoviesByGenres(genres []string, page int) (int, []*models.Movie, error) {
	startIndex := (page - 1) * _const.MoviesPageSize
	return moviesUC.movieRepository.GetMoviesByGenres(genres, startIndex)
}
