package movies

import "github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"

type UseCase interface {
	CreateMovie(movie *models.Movie) error

	GetMovie(id string) (*models.Movie, error)

	UpdateMovie(id string, newMovie *models.Movie) error
}