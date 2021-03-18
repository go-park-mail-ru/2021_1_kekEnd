package reviews

import "github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"

type ReviewRepository interface {
	CreateReview(review *models.Review) error

	CheckIfExists(username string, review *models.Review) bool

	GetUserReviews(username string) []*models.Review

	GetMovieReviews(movieID string) []*models.Review
}
