package usecase

import (
	"errors"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/reviews"
)

type ReviewsUseCase struct {
	reviewRepository reviews.ReviewRepository
}

func NewReviewsUseCase(repo reviews.ReviewRepository) *ReviewsUseCase {
	return &ReviewsUseCase{
		reviewRepository: repo,
	}
}

func (reviewsUC *ReviewsUseCase) CreateReview(username string, review *models.Review) error {
	_, err := reviewsUC.GetUserReviewForMovie(username, review.MovieID)
	if err == nil {
		return errors.New("review already exists")
	}
	return reviewsUC.reviewRepository.CreateReview(review)
}

func (reviewsUC *ReviewsUseCase) GetReviewsByUser(username string) []*models.Review {
	return reviewsUC.reviewRepository.GetUserReviews(username)
}

func (reviewsUC *ReviewsUseCase) GetReviewsByMovie(movieID string) []*models.Review {
	return reviewsUC.reviewRepository.GetMovieReviews(movieID)
}

func (reviewsUC *ReviewsUseCase) GetUserReviewForMovie(username string, movieID string) (*models.Review, error) {
	return reviewsUC.reviewRepository.GetUserReviewForMovie(username, movieID)
}
