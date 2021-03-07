package movies

import "github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"

type MovieRepository interface {
	CreateMovie(movie *models.Movie) error

	GetMovieById(id string) (*models.Movie, error)

	UpdateMovie(id string, newMovie *models.Movie) error
}