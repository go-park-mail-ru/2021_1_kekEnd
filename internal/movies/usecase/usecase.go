package usecase

import (
	"errors"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/movies"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/users"
	_const "github.com/go-park-mail-ru/2021_1_kekEnd/pkg/const"
)

type MoviesUseCase struct {
	movieRepository movies.MovieRepository
	userRepository  users.UserRepository
}

func NewMoviesUseCase(repo movies.MovieRepository, userRepo users.UserRepository) *MoviesUseCase {
	return &MoviesUseCase{
		movieRepository: repo,
		userRepository:  userRepo,
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

func (moviesUC *MoviesUseCase) MarkWatched(user models.User, id int) error {
	err := moviesUC.movieRepository.MarkWatched(user.Username, id)
	if err != nil {
		return err
	}
	// successful mark watched, must increment movies_watched for user
	newMoviesWatchNumber := *user.MoviesWatched + 1
	_, err = moviesUC.userRepository.UpdateUser(&user, models.User{
		Username:      user.Username,
		MoviesWatched: &newMoviesWatchNumber,
	})
	return err
}

func (moviesUC *MoviesUseCase) MarkUnwatched(user models.User, id int) error {
	err := moviesUC.movieRepository.MarkUnwatched(user.Username, id)
	if err != nil {
		return err
	}
	// successful mark unwatched, must decrement movies_watched for user
	newMoviesWatchNumber := *user.MoviesWatched - 1
	_, err = moviesUC.userRepository.UpdateUser(&user, models.User{
		Username:      user.Username,
		MoviesWatched: &newMoviesWatchNumber,
	})
	return err
}
