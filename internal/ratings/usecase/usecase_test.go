package usecase
//
//import (
//	"errors"
//	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
//	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/ratings/mocks"
//	"github.com/golang/mock/gomock"
//	"github.com/stretchr/testify/assert"
//	"testing"
//)
//
//func TestReviewsUseCase(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	repo := mocks.NewMockRepository(ctrl)
//	uc := NewRatingsUseCase(repo)
//
//	rating := &models.Rating{
//		UserID:  "let_robots_reign",
//		MovieID: "1",
//		Score:   7,
//	}
//
//	t.Run("CreateRating", func(t *testing.T) {
//
//	})
//
//	t.Run("GetRating", func(t *testing.T) {
//
//	})
//
//	t.Run("UpdateRating", func(t *testing.T) {
//
//	})
//
//	t.Run("DeleteRating", func(t *testing.T) {
//
//	})
//}