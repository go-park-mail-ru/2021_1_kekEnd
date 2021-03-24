package reviews

import "github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"

type ReviewRepository interface {
	CreateReview(review *models.Review) error

	GetUserReviews(username string) []*models.Review

	GetMovieReviews(movieID string) []*models.Review

	GetUserReviewForMovie(username string, movieID string) (*models.Review, error)

	DeleteUserReviewForMovie(username string, movieID string) error
}
