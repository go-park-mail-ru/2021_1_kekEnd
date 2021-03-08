package movies

import "github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"

type MovieRepository interface {
	CreateMovie(movie *models.Movie) error

	GetMovieByID(id string) (*models.Movie, error)
}
