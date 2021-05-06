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

	actor := models.Actor{
		ID:   "1",
		Name: "Tom Cruise",
	}

	t.Run("GetActor", func(t *testing.T) {
		repo.EXPECT().GetActorByID(actor.ID).Return(actor, nil)
		gotActor, err := uc.GetActor(actor.ID, actor.Name)
		assert.NoError(t, err)
		assert.Equal(t, actor, gotActor)
	})
}
