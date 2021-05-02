package movies

import "github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"

//go:generate mockgen -destination=mocks/usecase.go -package=mocks . UseCase
type UseCase interface {
	CreateMovie(movie *models.Movie) error

	GetMovie(id string, username string) (*models.Movie, error)

	GetBestMovies(page int, username string) (int, []*models.Movie, error)

	GetAllGenres() ([]string, error)

	GetMoviesByGenres(genres []string, page int, username string) (int, []*models.Movie, error)

	MarkWatched(username string, id int) error
}
