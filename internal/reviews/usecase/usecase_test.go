package usecase

import (
	"errors"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/reviews/mocks"
	userMocks "github.com/go-park-mail-ru/2021_1_kekEnd/internal/users/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReviewsUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockReviewRepository(ctrl)
	usersRepo := userMocks.NewMockUserRepository(ctrl)
	uc := NewReviewsUseCase(repo, usersRepo)

	var reviewsNum uint = 1
	user := &models.User{
		Username:      "let_robots_reign",
		Email:         "sample@ya.ru",
		Password:      "1234",
		ReviewsNumber: &reviewsNum,
	}

	review := &models.Review{
		ID:         "1",
		Title:      "Review",
		ReviewType: "positive",
		Content:    "test",
		MovieID:    "1",
	}

	t.Run("CreateReview", func(t *testing.T) {
		repo.EXPECT().GetUserReviewForMovie(user.Username, review.MovieID).Return(nil, errors.New("review doesn't exist"))
		repo.EXPECT().CreateReview(review).Return(nil)
		newReviewsNumber := *user.ReviewsNumber + 1
		usersRepo.EXPECT().UpdateUser(user, models.User{
			Username:      user.Username,
			ReviewsNumber: &newReviewsNumber,
		}).Return(user, nil)
		err := uc.CreateReview(user, review)
		assert.NoError(t, err)
	})

	t.Run("GetReviewsByUser", func(t *testing.T) {
		repo.EXPECT().GetUserReviews(user.Username).Return([]*models.Review{review}, nil)
		reviews, err := uc.GetReviewsByUser(user.Username)
		assert.NoError(t, err)
		assert.Equal(t, []*models.Review{review}, reviews)
	})

	t.Run("GetReviewsByMovie", func(t *testing.T) {
		repo.EXPECT().GetMovieReviews(review.MovieID, 0).Return(1, []*models.Review{review}, nil)
		pages, reviews, err := uc.GetReviewsByMovie(review.MovieID, 1)
		assert.NoError(t, err)
		assert.Equal(t, 1, pages)
		assert.Equal(t, []*models.Review{review}, reviews)
	})

	t.Run("GetUserReviewForMovie", func(t *testing.T) {
		repo.EXPECT().GetUserReviewForMovie(user.Username, review.MovieID).Return(review, nil)
		gotReview, err := uc.GetUserReviewForMovie(user.Username, review.MovieID)
		assert.NoError(t, err)
		assert.Equal(t, review, gotReview)
	})

	t.Run("EditUserReviewForMovie", func(t *testing.T) {
		newReview := &models.Review{
			ID:         "1",
			Title:      "New Review",
			ReviewType: "neutral",
			Content:    "content",
			MovieID:    "1",
		}
		repo.EXPECT().GetUserReviewForMovie(user.Username, review.MovieID).Return(review, nil)
		repo.EXPECT().EditUserReviewForMovie(newReview).Return(nil)
		err := uc.EditUserReviewForMovie(user, newReview)
		assert.NoError(t, err)
	})

	t.Run("DeleteUserReviewForMovie", func(t *testing.T) {
		repo.EXPECT().DeleteUserReviewForMovie(user.Username, review.MovieID).Return(nil)
		newReviewsNumber := *user.ReviewsNumber - 1
		usersRepo.EXPECT().UpdateUser(user, models.User{
			Username:      user.Username,
			ReviewsNumber: &newReviewsNumber,
		}).Return(user, nil)
		err := uc.DeleteUserReviewForMovie(user, review.MovieID)
		assert.NoError(t, err)
	})
}
