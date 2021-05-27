package usecase

import (
	"testing"

	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/actors/mocks"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestActorsUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepository(ctrl)
	uc := NewActorsUseCase(repo)

	user := &models.User{
		Username:      "let_robots_reign",
		Email:         "sample@mail.ru",
		Password:      "123",
		Avatar:        "http://localhost:8080/avatars/default.jpeg",
		MoviesWatched: new(uint),
		ReviewsNumber: new(uint),
		Subscribers:   new(uint),
		Subscriptions: new(uint),
	}

	actor := models.Actor{
		ID:   "1",
		Name: "Tom Cruise",
	}

	// t.Run("CreateActor", func(t *testing.T) {
	// 	repo.EXPECT().CreateActor(actor).Return(nil)
	// 	err := uc.CreateActor(*user, actor)
	// 	assert.NoError(t, err)
	// })

	t.Run("GetActor", func(t *testing.T) {
		repo.EXPECT().GetActorByID(actor.ID, "").Return(actor, nil)
		gotActor, err := uc.GetActor(actor.ID, "")
		assert.NoError(t, err)
		assert.Equal(t, actor, gotActor)
	})

	// t.Run("EditActor", func(t *testing.T) {
	// 	repo.EXPECT().EditActor(actor).Return(actor, nil)
	// 	gotActor, err := uc.EditActor(*user, actor)
	// 	assert.NoError(t, err)
	// 	assert.Equal(t, actor, gotActor)
	// })

	t.Run("LikeActor", func(t *testing.T) {
		repo.EXPECT().LikeActor(user.Username, 1).Return(nil)
		err := uc.LikeActor(user.Username, 1)
		assert.NoError(t, err)
	})

	t.Run("UnlikeActor", func(t *testing.T) {
		repo.EXPECT().UnlikeActor(user.Username, 1).Return(nil)
		err := uc.UnlikeActor(user.Username, 1)
		assert.NoError(t, err)
	})
}
