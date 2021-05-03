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
	_, err := moviesUC.movieRepository.GetMovieByID(movie.ID, "")
	if err == nil {
		return errors.New("movie already exists")
	}
	return moviesUC.movieRepository.CreateMovie(movie)
}

func (moviesUC *MoviesUseCase) GetMovie(id string, username string) (*models.Movie, error) {
	return moviesUC.movieRepository.GetMovieByID(id, username)
}

func (moviesUC *MoviesUseCase) GetBestMovies(page int, username string) (int, []*models.Movie, error) {
	startIndex := (page - 1) * _const.MoviesPageSize
	return moviesUC.movieRepository.GetBestMovies(startIndex, username)
}

func (moviesUC *MoviesUseCase) GetAllGenres() ([]string, error) {
	return moviesUC.movieRepository.GetAllGenres()
}

func (moviesUC *MoviesUseCase) GetMoviesByGenres(genres []string, page int, username string) (int, []*models.Movie, error) {
	startIndex := (page - 1) * _const.MoviesPageSize
	return moviesUC.movieRepository.GetMoviesByGenres(genres, startIndex, username)
}

func (moviesUC *MoviesUseCase) MarkWatched(username string, id int) error {
	return moviesUC.movieRepository.MarkWatched(username, id)
}

func (moviesUC *MoviesUseCase) MarkUnwatched(username string, id int) error {
	return moviesUC.movieRepository.MarkUnwatched(username, id)
}
