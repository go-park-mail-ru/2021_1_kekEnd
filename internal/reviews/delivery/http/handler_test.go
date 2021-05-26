package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/logger"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/middleware"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	reviewsMock "github.com/go-park-mail-ru/2021_1_kekEnd/internal/reviews/mocks"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/services/sessions/mocks"
	usersMock "github.com/go-park-mail-ru/2021_1_kekEnd/internal/users/mocks"
	"github.com/golang/mock/gomock"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestHandlers(t *testing.T) {
	r := gin.Default()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	lg := logger.NewAccessLogger()

	reviewsUC := reviewsMock.NewMockUseCase(ctrl)
	usersUC := usersMock.NewMockUseCase(ctrl)
	delivery := mocks.NewMockDelivery(ctrl)

	authMiddleware := middleware.NewAuthMiddleware(usersUC, delivery)
	api := r.Group("/api")
	v1 := api.Group("/v1")
	RegisterHttpEndpoints(v1, reviewsUC, usersUC, authMiddleware, lg)

	review := &models.Review{
		ID:         "1",
		Title:      "Review",
		ReviewType: "positive",
		Content:    "content",
		MovieID:    "1",
	}

	// wrongReview := &models.Review{
	// 	// ID:         "first",
	// 	Title:      "Review",
	// 	ReviewType: "pцositiveqwe",
	// 	Content:    "content",
	// 	MovieID:    "first",
	// }

	user := &models.User{
		Username:      "let_robots_reign",
		Email:         "sample@mail.ru",
		Password:      "123",
		Avatar:        "http://localhost:8080/avatars/default.jpeg",
		MoviesWatched: new(uint),
		ReviewsNumber: new(uint),
	}

	body, err := json.Marshal(review)
	assert.NoError(t, err)

	// wrongBody, err := json.Marshal(wrongReview)
	// assert.NoError(t, err)

	UUID := uuid.NewV4().String()

	t.Run("CreateReview", func(t *testing.T) {
		cookie := &http.Cookie{
			Name:  "session_id",
			Value: UUID,
		}

		usersUC.EXPECT().GetUser(user.Username).Return(user, nil).AnyTimes()

		delivery.EXPECT().GetUser(UUID).Return(user.Username, nil).AnyTimes()

		reviewsUC.EXPECT().CreateReview(user, review).Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/users/reviews", bytes.NewBuffer(body))
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("GetUserReviews", func(t *testing.T) {
		cookie := &http.Cookie{
			Name:  "session_id",
			Value: UUID,
		}

		usersUC.EXPECT().GetUser(user.Username).Return(user, nil).AnyTimes()

		delivery.EXPECT().GetUser(UUID).Return(user.Username, nil).AnyTimes()

		reviewsUC.EXPECT().GetReviewsByUser(user.Username).Return([]*models.Review{review}, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/user/%s/reviews", user.Username), nil)
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("GetMovieReviews", func(t *testing.T) {
		reviewsUC.EXPECT().GetReviewsByMovie(review.MovieID, 1).Return(1, []*models.Review{review}, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/movies/1/reviews", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("GetUserReviewForMovie", func(t *testing.T) {
		cookie := &http.Cookie{
			Name:  "session_id",
			Value: UUID,
		}

		usersUC.EXPECT().GetUser(user.Username).Return(user, nil).AnyTimes()

		delivery.EXPECT().GetUser(UUID).Return(user.Username, nil).AnyTimes()

		reviewsUC.EXPECT().GetUserReviewForMovie(user.Username, review.MovieID).Return(review, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/users/movies/1/reviews", nil)
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("EditUserReviewForMovie", func(t *testing.T) {
		cookie := &http.Cookie{
			Name:  "session_id",
			Value: UUID,
		}

		usersUC.EXPECT().GetUser(user.Username).Return(user, nil).AnyTimes()

		delivery.EXPECT().GetUser(UUID).Return(user.Username, nil).AnyTimes()

		newReview := &models.Review{
			ID:         "1",
			Title:      "New",
			ReviewType: "neutral",
			Content:    "new content",
			MovieID:    "1",
		}

		newBody, err := json.Marshal(newReview)
		assert.NoError(t, err)

		reviewsUC.EXPECT().EditUserReviewForMovie(user, newReview).Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/api/v1/users/movies/1/reviews", bytes.NewBuffer(newBody))
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("DeleteUserReviewForMovie", func(t *testing.T) {
		cookie := &http.Cookie{
			Name:  "session_id",
			Value: UUID,
		}

		usersUC.EXPECT().GetUser(user.Username).Return(user, nil).AnyTimes()

		delivery.EXPECT().GetUser(UUID).Return(user.Username, nil).AnyTimes()

		reviewsUC.EXPECT().DeleteUserReviewForMovie(user, review.MovieID).Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/v1/users/movies/1/reviews", nil)
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("CreateReviewError", func(t *testing.T) {
		cookie := &http.Cookie{
			Name:  "session_id",
			Value: UUID,
		}

		usersUC.EXPECT().GetUser(user.Username).Return(nil, errors.New("error")).AnyTimes()

		delivery.EXPECT().GetUser(UUID).Return("", errors.New("error")).AnyTimes()

		reviewsUC.EXPECT().CreateReview(user, review).Return(errors.New("error"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/users/reviews", bytes.NewBuffer(body))
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("GetUserReviewsError", func(t *testing.T) {
		cookie := &http.Cookie{
			Name:  "session_id",
			Value: UUID,
		}

		usersUC.EXPECT().GetUser("").Return(nil, errors.New("error")).AnyTimes()

		delivery.EXPECT().GetUser(UUID).Return(user.Username, nil).AnyTimes()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/user/%s/reviews", ""), nil)
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("GetMovieReviewsError", func(t *testing.T) {
		reviewsUC.EXPECT().GetReviewsByMovie(review.MovieID, 1).Return(1, nil, errors.New("error"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/movies/1/reviews", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("GetUserReviewForMovieError", func(t *testing.T) {
		cookie := &http.Cookie{
			Name:  "session_id",
			Value: UUID,
		}

		usersUC.EXPECT().GetUser(user.Username).Return(user, nil).AnyTimes()

		delivery.EXPECT().GetUser(UUID).Return(user.Username, nil).AnyTimes()

		reviewsUC.EXPECT().GetUserReviewForMovie(user.Username, review.MovieID).Return(review, errors.New("error"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/users/movies/1/reviews", nil)
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("EditUserReviewForMovieError", func(t *testing.T) {
		cookie := &http.Cookie{
			Name:  "session_id",
			Value: UUID,
		}

		usersUC.EXPECT().GetUser(user.Username).Return(user, nil).AnyTimes()

		delivery.EXPECT().GetUser(UUID).Return(user.Username, nil).AnyTimes()

		newReview := &models.Review{
			ID:         "qwe",
			Title:      "New",
			ReviewType: "neutral",
			Content:    "new content",
			MovieID:    "1",
		}

		newBody, err := json.Marshal(newReview)
		assert.NoError(t, err)

		reviewsUC.EXPECT().EditUserReviewForMovie(user, newReview).Return(errors.New("error"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/api/v1/users/movies/1/reviews", bytes.NewBuffer(newBody))
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("DeleteUserReviewForMovieError", func(t *testing.T) {
		cookie := &http.Cookie{
			Name:  "session_id",
			Value: UUID,
		}

		usersUC.EXPECT().GetUser(user.Username).Return(user, nil).AnyTimes()

		delivery.EXPECT().GetUser(UUID).Return(user.Username, nil).AnyTimes()

		reviewsUC.EXPECT().DeleteUserReviewForMovie(user, review.MovieID).Return(errors.New("error"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/v1/users/movies/1/reviews", nil)
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
