package usecase

import (
	"errors"

	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/ratings"
)

// RatingsUseCase структура usecase оценок
type RatingsUseCase struct {
	Repository ratings.Repository
}

// NewRatingsUseCase инициализация usecase оценок
func NewRatingsUseCase(repository ratings.Repository) *RatingsUseCase {
	return &RatingsUseCase{
		Repository: repository,
	}
}

// CreateRating создание оценки
func (u *RatingsUseCase) CreateRating(userID string, movieID string, score int) error {
	if score < 0 || score > 10 {
		return errors.New("invalid value for score")
	}

	return u.Repository.CreateRating(userID, movieID, score)
}

// GetRating получение оценки
func (u *RatingsUseCase) GetRating(userID string, movieID string) (models.Rating, error) {
	return u.Repository.GetRating(userID, movieID)
}

// UpdateRating обновление оценки
func (u *RatingsUseCase) UpdateRating(userID string, movieID string, score int) error {
	if score < 0 || score > 10 {
		return errors.New("invalid value for score")
	}

	return u.Repository.UpdateRating(userID, movieID, score)
}

// DeleteRating удаление оценки
func (u *RatingsUseCase) DeleteRating(userID string, movieID string) error {
	return u.Repository.DeleteRating(userID, movieID)
}
