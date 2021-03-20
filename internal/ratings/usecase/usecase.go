package usecase

import (
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/ratings"
)

type RatingUsecase struct {
	Repository ratings.Repository
}

func (u *RatingUsecase) CreateRating(userID string, movieID string, score uint) error {
	return u.Repository.CreateRating(userID, movieID, score)
}

func (u *RatingUsecase) GetRating(userID string, movieID string) (models.Rating, error) {
	return u.Repository.GetRating(userID, movieID)
}

func (u *RatingUsecase) DeleteRating(userID string, movieID string) error {
	return u.Repository.DeleteRating(userID, movieID)
}
