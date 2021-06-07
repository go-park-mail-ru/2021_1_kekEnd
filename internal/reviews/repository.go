package reviews

import "github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"

// ReviewRepository go:generate mockgen -destination=mocks/repository.go -package=mocks . ReviewRepository
type ReviewRepository interface {
	CreateReview(review *models.Review) error

	GetUserReviews(username string) ([]*models.Review, error)

	GetMovieReviews(movieID string, startInd int) (int, []*models.Review, error)

	GetUserReviewForMovie(username string, movieID string) (*models.Review, error)

	EditUserReviewForMovie(review *models.Review) error

	DeleteUserReviewForMovie(username string, movieID string) error

	GetFeed([]models.UserNoPassword) ([]models.ReviewFeedItem, error)

	DeleteReview(username string, movieID int) error
}
