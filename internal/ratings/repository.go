package ratings

import "github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"

type Repository interface {
	CreateRating(username string, movieID string, score int) error
	GetRating(username string, movieID string) (models.Rating, error)
	UpdateRating(username string, movieID string, score int) error
	DeleteRating(username string, movieID string) error
}
