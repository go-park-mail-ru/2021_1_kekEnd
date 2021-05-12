package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	uuid "github.com/satori/go.uuid"

	"github.com/gin-gonic/gin"
	actorsMock "github.com/go-park-mail-ru/2021_1_kekEnd/internal/actors/mocks"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/logger"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/middleware"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/services/sessions/mocks"
	usersMock "github.com/go-park-mail-ru/2021_1_kekEnd/internal/users/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHandlers(t *testing.T) {
	r := gin.Default()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	lg := logger.NewAccessLogger()

	actorsUC := actorsMock.NewMockUseCase(ctrl)
	usersUC := usersMock.NewMockUseCase(ctrl)
	delivery := mocks.NewMockDelivery(ctrl)

	authMiddleware := middleware.NewAuthMiddleware(usersUC, delivery)

	RegisterHttpEndpoints(r, actorsUC, authMiddleware, lg)

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
		ID:          "1",
		Name:        "Tom Cruise",
		MoviesCount: 1,
	}
	newBody, err := json.Marshal(actor)
	assert.NoError(t, err)

	usersUC.EXPECT().GetUser(user.Username).Return(user, nil).AnyTimes()
	UUID := uuid.NewV4().String()
	delivery.EXPECT().GetUser(UUID).Return(user.Username, nil).AnyTimes()

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: UUID,
	}

	t.Run("CreateActor", func(t *testing.T) {
		actorsUC.EXPECT().CreateActor(*user, actor).Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/actors", bytes.NewBuffer(newBody))
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("GetActor", func(t *testing.T) {
		actorsUC.EXPECT().GetActor(actor.ID, "").Return(actor, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/actors/1", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("EditActor", func(t *testing.T) {
		actorsUC.EXPECT().EditActor(*user, actor).Return(actor, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/actors/1", bytes.NewBuffer(newBody))
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("LikeActor", func(t *testing.T) {
		actorsUC.EXPECT().LikeActor(user.Username, 1).Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/actors/1/like", nil)
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("UnlikeActor", func(t *testing.T) {
		actorsUC.EXPECT().UnlikeActor(user.Username, 1).Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/actors/1/like", nil)
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("CreateActorError", func(t *testing.T) {
		actorsUC.EXPECT().CreateActor(*user, actor).Return(errors.New("error"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/actors", bytes.NewBuffer(newBody))
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("CreateActorError2", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/actors", nil)
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("GetActorError", func(t *testing.T) {
		actorsUC.EXPECT().GetActor(actor.ID, "").Return(models.Actor{}, errors.New("error"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/actors/1", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("EditActorError", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/actors/1", nil)
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("EditActorError2", func(t *testing.T) {
		actorsUC.EXPECT().EditActor(*user, actor).Return(models.Actor{}, errors.New("error"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/actors/1", bytes.NewBuffer(newBody))
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("LikeActorError", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/actors/:actor_id/like", nil)
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("LikeActorError2", func(t *testing.T) {
		actorsUC.EXPECT().LikeActor(user.Username, 1).Return(errors.New("error"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/actors/1/like", nil)
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("UnlikeActorError", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/actors/:actor_id/like", nil)
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("UnlikeActorError2", func(t *testing.T) {
		actorsUC.EXPECT().UnlikeActor(user.Username, 1).Return(errors.New("error"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/actors/1/like", nil)
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
