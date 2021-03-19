package reviews

import "github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"

type UseCase interface {
	CreateReview(username string, review *models.Review) error

	GetReviewsByUser(username string) []*models.Review

	GetReviewsByMovie(movieID string) []*models.Review

	GetUserReviewForMovie(username string, movieID string) (*models.Review, error)
}
