package sessions

import (
	"errors"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/services/sessions/mocks"
	"github.com/golang/mock/gomock"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	testErr := errors.New("error in Create Session")

	t.Run("Create-OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		rdb := mocks.NewMockRepository(ctrl)
		useCase := NewUseCase(rdb)

		username := "whaevaforeva"

		rdb.
			EXPECT().
			Create(gomock.Any(), username, time.Duration(10)).
			Return(nil)

		_, err := useCase.Create(username, time.Duration(10))
		assert.NoError(t, err)
	})

	t.Run("Create-Error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		rdb := mocks.NewMockRepository(ctrl)
		useCase := NewUseCase(rdb)

		username := "whaevaforeva"

		rdb.
			EXPECT().
			Create(gomock.Any(), username, time.Duration(10)).
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

		rdb := mocks.NewMockRepository(ctrl)
		useCase := NewUseCase(rdb)

		username := "whaevaforeva"
		ID := uuid.NewV4().String()
		sessionID := addPrefix(ID)

		rdb.
			EXPECT().
			Get(sessionID).
			Return(username, nil)

		userFromSession, err := useCase.GetUser(ID)

		assert.NoError(t, err)
		assert.Equal(t, username, userFromSession)
	})

	t.Run("GetUser-Error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		rdb := mocks.NewMockRepository(ctrl)
		useCase := NewUseCase(rdb)

		username := "whaevaforeva"
		ID := uuid.NewV4().String()
		sessionID := addPrefix(ID)

		rdb.
			EXPECT().
			Get(sessionID).
			Return(username, testErr)

		_, err := useCase.GetUser(ID)
		assert.Error(t, err)
	})
}

func TestDelete(t *testing.T) {
	testErr := errors.New("error in Delete Session")

	t.Run("Delete-OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		rdb := mocks.NewMockRepository(ctrl)
		useCase := NewUseCase(rdb)

		ID := uuid.NewV4().String()
		sessionID := addPrefix(ID)
		rdb.
			EXPECT().
			Delete(sessionID).
			Return(nil)

		err := useCase.Delete(ID)

		assert.NoError(t, err)
	})

	t.Run("Delete-Error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		rdb := mocks.NewMockRepository(ctrl)
		useCase := NewUseCase(rdb)

		ID := uuid.NewV4().String()
		sessionID := addPrefix(ID)
		rdb.
			EXPECT().
			Delete(sessionID).
			Return(testErr)

		err := useCase.Delete(ID)
		assert.Error(t, err)
	})
}
