package usecase

import (
	"errors"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/ratings"
)

type RatingsUseCase struct {
	Repository ratings.Repository
}

func NewRatingsUseCase(repository ratings.Repository) *RatingsUseCase {
	return &RatingsUseCase{
		Repository: repository,
	}
}

func (u *RatingsUseCase) CreateRating(userID string, movieID string, score int) error {
	if score < 0 || score > 10 {
		return errors.New("invalid value for score")
	}

	return u.Repository.CreateRating(userID, movieID, score)
}

func (u *RatingsUseCase) GetRating(userID string, movieID string) (models.Rating, error) {
	return u.Repository.GetRating(userID, movieID)
}

func (u *RatingsUseCase) UpdateRating(userID string, movieID string, score int) error {
	if score < 0 || score > 10 {
		return errors.New("invalid value for score")
	}

	return u.Repository.UpdateRating(userID, movieID, score)
}

func (u *RatingsUseCase) DeleteRating(userID string, movieID string) error {
	return u.Repository.DeleteRating(userID, movieID)
}
