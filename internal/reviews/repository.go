package reviews

import "github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"

type ReviewRepository interface {
	CreateReview(review *models.Review) error

	GetUserReviews(username string) ([]*models.Review, error)

	GetMovieReviews(movieID string, startInd int) (int, []*models.Review, error)

	GetUserReviewForMovie(username string, movieID string) (*models.Review, error)

	EditUserReviewForMovie(review *models.Review) error

	DeleteUserReviewForMovie(username string, movieID string) error
}
