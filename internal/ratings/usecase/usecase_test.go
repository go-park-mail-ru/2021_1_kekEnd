package usecase

import (
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/ratings/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReviewsUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepository(ctrl)
	uc := NewRatingsUseCase(repo)

	rating := &models.Rating{
		UserID:  "let_robots_reign",
		MovieID: "1",
		Score:   7,
	}

	t.Run("CreateRating", func(t *testing.T) {
		repo.EXPECT().CreateRating(rating.UserID, rating.MovieID, rating.Score).Return(nil)
		err := uc.CreateRating(rating.UserID, rating.MovieID, rating.Score)
		assert.NoError(t, err)
	})

	t.Run("GetRating", func(t *testing.T) {
		repo.EXPECT().GetRating(rating.UserID, rating.MovieID).Return(*rating, nil)
		gotRating, err := uc.GetRating(rating.UserID, rating.MovieID)
		assert.NoError(t, err)
		assert.Equal(t, *rating, gotRating)
	})

	t.Run("UpdateRating", func(t *testing.T) {
		newRating := &models.Rating{
			UserID:  "1",
			MovieID: "1",
			Score:   10,
		}
		repo.EXPECT().UpdateRating(newRating.UserID, newRating.MovieID, newRating.Score).Return(nil)
		err := uc.UpdateRating(newRating.UserID, newRating.MovieID, newRating.Score)
		assert.NoError(t, err)
	})

	t.Run("DeleteRating", func(t *testing.T) {
		repo.EXPECT().DeleteRating(rating.UserID, rating.MovieID).Return(nil)
		err := uc.DeleteRating(rating.UserID, rating.MovieID)
		assert.NoError(t, err)
	})
}
