package reviews

import "github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"

type UseCase interface {
	CreateReview(user *models.User, review *models.Review) error

	GetReviewsByUser(username string) []*models.Review

	GetReviewsByMovie(movieID string) []*models.Review

	GetUserReviewForMovie(username string, movieID string) (*models.Review, error)

	DeleteUserReviewForMovie(username string, movieID string) error
}
