package sessions

import (
	"errors"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/sessions"
	"github.com/golang/mock/gomock"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	testErr := errors.New("error in Create Session")

	t.Run("Create-OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		rdb := sessions.NewMockRepository(ctrl)
		useCase := NewUseCase(rdb)

		username := "whaevaforeva"
		//UUID := addPrefix(uuid.NewV4().String())

		rdb.
			EXPECT().
			Create(mock.Anything, username, time.Duration(10)).
			Return(nil)

		_, err := useCase.Create(username, time.Duration(10))
		assert.NoError(t, err)
	})

	t.Run("Create-Error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		rdb := sessions.NewMockRepository(ctrl)
		useCase := NewUseCase(rdb)

		username := "whaevaforeva"

		rdb.
			EXPECT().
			Create(mock.Anything, username, time.Duration(10)).
			Return(testErr)

		_, err := useCase.Create(username, time.Duration(10))
		assert.Error(t, err)
	})
}

func TestGetUser(t *testing.T) {
	testErr := errors.New("error in GetUser from Session")

	t.Run("GetUser-OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		rdb := sessions.NewMockRepository(ctrl)
		useCase := NewUseCase(rdb)

		username := "whaevaforeva"
		sessionID := addPrefix(uuid.NewV4().String())

		rdb.
			EXPECT().
			Get(sessionID).
			Return(username, nil)

		userFromSession, err := useCase.Check(sessionID)

		assert.NoError(t, err)
		assert.Equal(t, username, userFromSession)
	})

	t.Run("GetUser-Error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		rdb := sessions.NewMockRepository(ctrl)
		useCase := NewUseCase(rdb)

		username := "whaevaforeva"
		sessionID := addPrefix(uuid.NewV4().String())

		rdb.
			EXPECT().
			Get(sessionID).
			Return(username, testErr)

		_, err := useCase.Check(sessionID)
		assert.Error(t, err)
	})
}

func TestDelete(t *testing.T) {
	testErr := errors.New("error in Delete Session")

	t.Run("Delete-OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		rdb := sessions.NewMockRepository(ctrl)
		useCase := NewUseCase(rdb)

		sessionID := addPrefix(uuid.NewV4().String())

		rdb.
			EXPECT().
			Delete(sessionID).
			Return(nil)

		err := useCase.Delete(sessionID)

		assert.NoError(t, err)
	})

	t.Run("Delete-Error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		rdb := sessions.NewMockRepository(ctrl)
		useCase := NewUseCase(rdb)
		sessionID := addPrefix(uuid.NewV4().String())

		rdb.
			EXPECT().
			Delete(sessionID).
			Return(testErr)

		err := useCase.Delete(sessionID)
		assert.Error(t, err)
	})
}