package http

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/middleware"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	reviewsMock "github.com/go-park-mail-ru/2021_1_kekEnd/internal/reviews/mocks"
	sessionsMock "github.com/go-park-mail-ru/2021_1_kekEnd/internal/sessions"
	sessions "github.com/go-park-mail-ru/2021_1_kekEnd/internal/sessions/delivery"
	usersMock "github.com/go-park-mail-ru/2021_1_kekEnd/internal/users/mocks"
	"github.com/golang/mock/gomock"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlers(t *testing.T) {
	r := gin.Default()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reviewsUC := reviewsMock.NewMockUseCase(ctrl)
	usersUC := usersMock.NewMockUseCase(ctrl)
	sessionsUC := sessionsMock.NewMockUseCase(ctrl)
	delivery := sessions.NewDelivery(sessionsUC)

	authMiddleware := middleware.NewAuthMiddleware(usersUC, delivery)

	RegisterHttpEndpoints(r, reviewsUC, usersUC, authMiddleware)

	createBody := &models.Review{
		ID:         "1",
		Title:      "Review",
		ReviewType: "positive",
		Content:    "content",
		MovieID:    "1",
	}

	user := &models.User{
		Username:      "let_robots_reign",
		Email:         "sample@mail.ru",
		Password:      "123",
		Avatar:        "http://localhost:8080/avatars/default.jpeg",
		MoviesWatched: new(uint),
		ReviewsNumber: new(uint),
	}

	body, err := json.Marshal(createBody)
	assert.NoError(t, err)

	UUID := uuid.NewV4().String()

	t.Run("CreateReview", func(t *testing.T) {
		cookie := &http.Cookie{
			Name:  "session_id",
			Value: UUID,
		}

		usersUC.EXPECT().GetUser(user.Username).Return(user, nil).AnyTimes()

		sessionsUC.
			EXPECT().
			Check(UUID).
			Return(user.Username, nil).AnyTimes()

		reviewsUC.EXPECT().CreateReview(user, createBody).Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/users/reviews", bytes.NewBuffer(body))
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})
}