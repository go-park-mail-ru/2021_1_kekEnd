package usecase

import (
	"errors"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/reviews"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/users"
)

type ReviewsUseCase struct {
	reviewRepository reviews.ReviewRepository
	userRepository   users.UserRepository
}

func NewReviewsUseCase(reviewRepo reviews.ReviewRepository, userRepo users.UserRepository) *ReviewsUseCase {
	return &ReviewsUseCase{
		reviewRepository: reviewRepo,
		userRepository:   userRepo,
	}
}

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
	_, err = reviewsUC.userRepository.UpdateUser(user, models.User{
		Username:      user.Username,
		ReviewsNumber: user.ReviewsNumber + 1,
	})
	return err
}

func (reviewsUC *ReviewsUseCase) GetReviewsByUser(username string) []*models.Review {
	return reviewsUC.reviewRepository.GetUserReviews(username)
}

func (reviewsUC *ReviewsUseCase) GetReviewsByMovie(movieID string) []*models.Review {
	return reviewsUC.reviewRepository.GetMovieReviews(movieID)
}

func (reviewsUC *ReviewsUseCase) GetUserReviewForMovie(username string, movieID string) ([]*models.Review, error) {
	return reviewsUC.reviewRepository.GetUserReviewForMovie(username, movieID)
}

func (reviewsUC *ReviewsUseCase) EditUserReviewForMovie(user *models.User, review *models.Review) error {
	_, err := reviewsUC.GetUserReviewForMovie(user.Username, review.MovieID)
	if err != nil {
		return errors.New("review doesn't exist")
	}
	return reviewsUC.reviewRepository.EditUserReviewForMovie(review)
}

func (reviewsUC *ReviewsUseCase) DeleteUserReviewForMovie(user *models.User, movieID string) error {
	err := reviewsUC.reviewRepository.DeleteUserReviewForMovie(user.Username, movieID)
	if err != nil {
		return err
	}
	// successful deletion, must decrement reviews_number for user
	// TODO: when decrementing from 1 to 0, it doesn't update because of 'if change.ReviewsNumber != 0' in repository
	_, err = reviewsUC.userRepository.UpdateUser(user, models.User{
		Username:      user.Username,
		ReviewsNumber: user.ReviewsNumber - 1,
	})
	return err
}
