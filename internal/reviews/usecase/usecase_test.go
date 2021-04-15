package usecase

import (
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/reviews/mocks"
	"github.com/golang/mock/gomock"
	"testing"
)

func TestReviewsUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	
	repo := mocks.NewMockReviewRepository(ctrl)
	uc := NewReviewsUseCase(repo)
	
	review := &models.Review{
		ID:         "1",
		Title:      "Review",
		ReviewType: "positive",
		Content:    "test",
		Author:     "let_robots_reign",
		MovieID:    "1",
	}
}