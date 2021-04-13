package movies

import "github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"

type UseCase interface {
	CreateMovie(movie *models.Movie) error

	GetMovie(id string) (*models.Movie, error)

	GetBestMovies(page int) (int, []*models.Movie, error)

	GetAllGenres() ([]string, error)

	GetMoviesByGenres(genres []string, page int) (int, []*models.Movie, error)
}
