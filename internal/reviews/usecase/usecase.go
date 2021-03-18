package usecase

import (
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

func (reviewsUC *ReviewsUseCase) CreateReview(review *models.Review) error {
	// TODO: check if user already uploaded a review for this movie
	return reviewsUC.reviewRepository.CreateReview(review)
}

func (reviewsUC *ReviewsUseCase) GetReviewsByUser(username string) []*models.Review {
	return reviewsUC.reviewRepository.GetUserReviews(username)
}

func (reviewsUC *ReviewsUseCase) GetReviewsByMovie(movieID string) []*models.Review {
	return reviewsUC.reviewRepository.GetMovieReviews(movieID)
}
