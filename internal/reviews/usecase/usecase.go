package usecase

import (
	"errors"

	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/reviews"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/users"
	constants "github.com/go-park-mail-ru/2021_1_kekEnd/pkg/const"
)

// ReviewsUseCase структура usecase рецензий
type ReviewsUseCase struct {
	reviewRepository reviews.ReviewRepository
	userRepository   users.UserRepository
}

// NewReviewsUseCase инициализация структуры usecase рецензий
func NewReviewsUseCase(reviewRepo reviews.ReviewRepository, userRepo users.UserRepository) *ReviewsUseCase {
	return &ReviewsUseCase{
		reviewRepository: reviewRepo,
		userRepository:   userRepo,
	}
}

// CreateReview создание рецензии
func (reviewsUC *ReviewsUseCase) CreateReview(user *models.User, review *models.Review) error {
	_, err := reviewsUC.GetUserReviewForMovie(user.Username, review.MovieID)
	if err == nil {
		return errors.New("review already exists")
	}

	review.Author = user.Username
	err = reviewsUC.reviewRepository.CreateReview(review)
	if err != nil {
		return err
	}
	// successful create, must increment reviews_number for user
	newReviewsNumber := *user.ReviewsNumber + 1
	_, err = reviewsUC.userRepository.UpdateUser(user, models.User{
		Username:      user.Username,
		ReviewsNumber: &newReviewsNumber,
	})
	return err
}

// GetReviewsByUser полученить рецензии пользователя
func (reviewsUC *ReviewsUseCase) GetReviewsByUser(username string) ([]*models.Review, error) {
	return reviewsUC.reviewRepository.GetUserReviews(username)
}

// GetReviewsByMovie получить рецензии фильма
func (reviewsUC *ReviewsUseCase) GetReviewsByMovie(movieID string, page int) (int, []*models.Review, error) {
	startIndex := (page - 1) * constants.ReviewsPageSize

	return reviewsUC.reviewRepository.GetMovieReviews(movieID, startIndex)
}

// GetUserReviewForMovie получить рецензию пользователя к фильму
func (reviewsUC *ReviewsUseCase) GetUserReviewForMovie(username string, movieID string) (*models.Review, error) {
	return reviewsUC.reviewRepository.GetUserReviewForMovie(username, movieID)
}

// EditUserReviewForMovie изменить рецензию пользователя
func (reviewsUC *ReviewsUseCase) EditUserReviewForMovie(user *models.User, review *models.Review) error {
	oldReview, err := reviewsUC.GetUserReviewForMovie(user.Username, review.MovieID)
	if err != nil {
		return errors.New("review doesn't exist")
	}
	review.ID = oldReview.ID
	review.Author = user.Username
	return reviewsUC.reviewRepository.EditUserReviewForMovie(review)
}

// DeleteUserReviewForMovie удалить рецензнию пользователя
func (reviewsUC *ReviewsUseCase) DeleteUserReviewForMovie(user *models.User, movieID string) error {
	err := reviewsUC.reviewRepository.DeleteUserReviewForMovie(user.Username, movieID)
	if err != nil {
		return err
	}
	// successful deletion, must decrement reviews_number for user
	newReviewsNumber := *user.ReviewsNumber - 1
	_, err = reviewsUC.userRepository.UpdateUser(user, models.User{
		Username:      user.Username,
		ReviewsNumber: &newReviewsNumber,
	})
	return err
}

// DeleteUserReviewForMovie удалить рецензнию пользователя
func (reviewsUC *ReviewsUseCase) DeleteReview(admin string, username string, movieID int) error {
	if admin != "admin1" {
		return errors.New("dont have permission")
	}

	return reviewsUC.reviewRepository.DeleteReview(username, movieID)
	// // successful deletion, must decrement reviews_number for user
	// newReviewsNumber := *user.ReviewsNumber - 1
	// _, err = reviewsUC.userRepository.UpdateUser(user, models.User{
	// 	Username:      username,
	// 	ReviewsNumber: &newReviewsNumber,
	// })
	// return err
}
