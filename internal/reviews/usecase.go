package reviews

import "github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"

type UseCase interface {
	CreateReview(user *models.User, review *models.Review) error

	GetReviewsByUser(username string) ([]*models.Review, error)

	GetReviewsByMovie(movieID string, page int) (int, []*models.Review, error)

	GetUserReviewForMovie(username string, movieID string) (*models.Review, error)

	EditUserReviewForMovie(user *models.User, review *models.Review) error

	DeleteUserReviewForMovie(user *models.User, movieID string) error
}
