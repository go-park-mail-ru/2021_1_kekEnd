package sessions

import (
	"context"
	"errors"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/sessions"
	"github.com/golang/mock/gomock"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreate(t *testing.T) {
	testErr := errors.New("error in Create Session")

	t.Run("Create-OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mUC := sessions.NewMockUseCase(ctrl)
		delivery := NewDelivery(mUC)

		username := "whaevaforeva"
		UUID := uuid.NewV4().String()

		mUC.
			EXPECT().
			Create(context.TODO(), username, 123123).
			Return(UUID, nil)

		sessionID, err := delivery.Create(context.TODO(), username, 123123)
		assert.NoError(t, err)
		assert.Equal(t, UUID, sessionID)
	})

	t.Run("Create-Error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mUC := sessions.NewMockUseCase(ctrl)
		delivery := NewDelivery(mUC)

		username := "whaevaforeva"
		UUID := uuid.NewV4().String()

		mUC.
			EXPECT().
			Create(context.TODO(), username, 123123).
			Return(UUID, testErr)

		_, err := delivery.Create(context.TODO(), username, 12323)
		assert.Error(t, err)
	})
}
