package movies

import "github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"

//go:generate mockgen -destination=mocks/repository.go -package=mocks . MovieRepository
type MovieRepository interface {
	CreateMovie(movie *models.Movie) error

	GetMovieByID(id string, username string) (*models.Movie, error)

	GetBestMovies(startIndex int, username string) (int, []*models.Movie, error)

	GetAllGenres() ([]string, error)

	GetMoviesByGenres(genres []string, startIndex int, username string) (int, []*models.Movie, error)

	MarkWatched(username string, id int) error
}
