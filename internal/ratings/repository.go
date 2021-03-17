package ratings

import "github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"

type Repository interface {
	CreateUserRating(userID string, movieID string, score uint) error
	GetUserRating(userID string, movieID string) models.Rating
	DeleteUserRating(userID string, movieID string) error
}
