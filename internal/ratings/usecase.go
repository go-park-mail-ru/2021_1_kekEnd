package ratings

import "github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"

type UseCase interface {
	CreateRating(userID string, movieID string, score uint) error
	GetRating(userID string, movieID string) (models.Rating, error)
	DeleteRating(userID string, movieID string) error
}
