package sessions

import (
	"errors"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/sessions"
	"github.com/golang/mock/gomock"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
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
			Create(username, time.Duration(10)).
			Return(UUID, nil)

		sessionID, err := delivery.Create(username, time.Duration(10))
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
			Create(username, time.Duration(10)).
			Return(UUID, testErr)

		_, err := delivery.Create(username, time.Duration(10))
		assert.Error(t, err)
	})
}

func TestGetUser(t *testing.T) {
	testErr := errors.New("error in GetUser from Session")

	t.Run("GetUser-OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mUC := sessions.NewMockUseCase(ctrl)
		delivery := NewDelivery(mUC)

		username := "whaevaforeva"
		sessionID := uuid.NewV4().String()

		mUC.
			EXPECT().
			Check(sessionID).
			Return(username, nil)

		userFromSession, err := delivery.GetUser(sessionID)

		assert.NoError(t, err)
		assert.Equal(t, username, userFromSession)
	})

	t.Run("GetUser-Error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mUC := sessions.NewMockUseCase(ctrl)
		delivery := NewDelivery(mUC)

		username := "whaevaforeva"
		sessionID := uuid.NewV4().String()

		mUC.
			EXPECT().
			Check(sessionID).
			Return(username, testErr)

		_, err := delivery.GetUser(sessionID)

		assert.Error(t, err)
	})
}

func TestDelete(t *testing.T) {
	testErr := errors.New("error in Delete Session")

	t.Run("Delete-OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mUC := sessions.NewMockUseCase(ctrl)
		delivery := NewDelivery(mUC)

		sessionID := uuid.NewV4().String()

		mUC.
			EXPECT().
			Delete(sessionID).
			Return(nil)

		err := delivery.Delete(sessionID)

		assert.NoError(t, err)
	})

	t.Run("Delete-Error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mUC := sessions.NewMockUseCase(ctrl)
		delivery := NewDelivery(mUC)
		sessionID := uuid.NewV4().String()

		mUC.
			EXPECT().
			Delete(sessionID).
			Return(testErr)

		_, err := delivery.GetUser(sessionID)

		assert.Error(t, err)
	})
}