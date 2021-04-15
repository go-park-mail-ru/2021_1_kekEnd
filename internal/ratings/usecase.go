package ratings

import "github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"

//go:generate mockgen -destination=mocks/usecase.go -package=mocks . UseCase
type UseCase interface {
	CreateRating(userID string, movieID string, score int) error
	GetRating(userID string, movieID string) (models.Rating, error)
	UpdateRating(userID string, movieID string, score int) error
	DeleteRating(userID string, movieID string) error
}
