package ratings

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/logger"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/middleware"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	ratingsMock "github.com/go-park-mail-ru/2021_1_kekEnd/internal/ratings/mocks"
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

	ratingsUC := ratingsMock.NewMockUseCase(ctrl)
	usersUC := usersMock.NewMockUseCase(ctrl)
	delivery := mocks.NewMockDelivery(ctrl)

	authMiddleware := middleware.NewAuthMiddleware(usersUC, delivery)
	api := r.Group("/api")
	v1 := api.Group("/v1")
	RegisterHTTPEndpoints(v1, ratingsUC, authMiddleware, lg)

	data := ratingData{
		MovieID: "1",
		Score:   "8",
	}
	wrongData := ratingData{
		MovieID: "str1",
		Score:   "str2",
	}
	rating := models.Rating{
		UserID:  "let_robots_reign",
		MovieID: "1",
		Score:   8,
	}

	user := &models.User{
		Username:      "let_robots_reign",
		Email:         "sample@mail.ru",
		Password:      "123",
		Avatar:        "http://localhost:8080/avatars/default.jpeg",
		MoviesWatched: new(uint),
		ReviewsNumber: new(uint),
	}

	body, err := json.Marshal(data)
	assert.NoError(t, err)

	wrongBody, err := json.Marshal(wrongData)
	assert.NoError(t, err)

	UUID := uuid.NewV4().String()

	t.Run("CreateRating", func(t *testing.T) {
		cookie := &http.Cookie{
			Name:  "session_id",
			Value: UUID,
		}

		usersUC.EXPECT().GetUser(user.Username).Return(user, nil).AnyTimes()

		delivery.EXPECT().GetUser(UUID).Return(user.Username, nil).AnyTimes()

		ratingsUC.EXPECT().CreateRating(user.Username, rating.MovieID, rating.Score).Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/ratings", bytes.NewBuffer(body))
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("GetRating", func(t *testing.T) {
		cookie := &http.Cookie{
			Name:  "session_id",
			Value: UUID,
		}

		usersUC.EXPECT().GetUser(user.Username).Return(user, nil).AnyTimes()

		delivery.EXPECT().GetUser(UUID).Return(user.Username, nil).AnyTimes()

		ratingsUC.EXPECT().GetRating(user.Username, rating.MovieID).Return(rating, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/ratings/1", nil)
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("UpdateRating", func(t *testing.T) {
		cookie := &http.Cookie{
			Name:  "session_id",
			Value: UUID,
		}

		usersUC.EXPECT().GetUser(user.Username).Return(user, nil).AnyTimes()

		delivery.EXPECT().GetUser(UUID).Return(user.Username, nil).AnyTimes()

		newRating := models.Rating{
			MovieID: "1",
			Score:   10,
		}

		newData := ratingData{
			MovieID: "1",
			Score:   "10",
		}

		newBody, err := json.Marshal(newData)
		assert.NoError(t, err)

		ratingsUC.EXPECT().UpdateRating(user.Username, newRating.MovieID, newRating.Score).Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/api/v1/ratings", bytes.NewBuffer(newBody))
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("DeleteRating", func(t *testing.T) {
		cookie := &http.Cookie{
			Name:  "session_id",
			Value: UUID,
		}

		usersUC.EXPECT().GetUser(user.Username).Return(user, nil).AnyTimes()

		delivery.EXPECT().GetUser(UUID).Return(user.Username, nil).AnyTimes()

		ratingsUC.EXPECT().DeleteRating(user.Username, rating.MovieID).Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/v1/ratings/1", nil)
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("CreateRatingError", func(t *testing.T) {
		cookie := &http.Cookie{
			Name:  "session_id",
			Value: UUID,
		}

		usersUC.EXPECT().GetUser(user.Username).Return(nil, errors.New("error")).AnyTimes()

		delivery.EXPECT().GetUser(UUID).Return("", errors.New("error")).AnyTimes()

		ratingsUC.EXPECT().CreateRating(user.Username, rating.MovieID, rating.Score).Return(errors.New("error"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/ratings", bytes.NewBuffer(body))
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("CreateRatingError2", func(t *testing.T) {
		cookie := &http.Cookie{
			Name:  "session_id",
			Value: UUID,
		}

		usersUC.EXPECT().GetUser(user.Username).Return(nil, errors.New("error")).AnyTimes()

		delivery.EXPECT().GetUser(UUID).Return("", errors.New("error")).AnyTimes()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/ratings", bytes.NewBuffer(wrongBody))
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("GetRatingError", func(t *testing.T) {
		cookie := &http.Cookie{
			Name:  "session_id",
			Value: UUID,
		}

		usersUC.EXPECT().GetUser(user.Username).Return(user, errors.New("error")).AnyTimes()

		delivery.EXPECT().GetUser(UUID).Return(user.Username, errors.New("error")).AnyTimes()

		ratingsUC.EXPECT().GetRating(user.Username, rating.MovieID).Return(rating, errors.New("error"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/ratings/1", nil)
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("UpdateRatingError", func(t *testing.T) {
		cookie := &http.Cookie{
			Name:  "session_id",
			Value: UUID,
		}

		usersUC.EXPECT().GetUser(user.Username).Return(user, errors.New("error")).AnyTimes()

		delivery.EXPECT().GetUser(UUID).Return(user.Username, errors.New("error")).AnyTimes()

		newRating := models.Rating{
			MovieID: "1",
			Score:   10,
		}

		newData := ratingData{
			MovieID: "1",
			Score:   "10",
		}

		newBody, err := json.Marshal(newData)
		assert.NoError(t, err)

		ratingsUC.EXPECT().UpdateRating(user.Username, newRating.MovieID, newRating.Score).Return(errors.New("error"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/api/v1/ratings", bytes.NewBuffer(newBody))
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("UpdateRatingError2", func(t *testing.T) {
		cookie := &http.Cookie{
			Name:  "session_id",
			Value: UUID,
		}

		usersUC.EXPECT().GetUser(user.Username).Return(user, errors.New("error")).AnyTimes()

		delivery.EXPECT().GetUser(UUID).Return(user.Username, errors.New("error")).AnyTimes()

		newData := ratingData{
			MovieID: "qwe",
			Score:   "qwe",
		}

		newBody, err := json.Marshal(newData)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/api/v1/ratings", bytes.NewBuffer(newBody))
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("DeleteRatingError", func(t *testing.T) {
		cookie := &http.Cookie{
			Name:  "session_id",
			Value: UUID,
		}

		usersUC.EXPECT().GetUser(user.Username).Return(nil, errors.New("error")).AnyTimes()

		delivery.EXPECT().GetUser(UUID).Return(user.Username, errors.New("error")).AnyTimes()

		ratingsUC.EXPECT().DeleteRating(user.Username, rating.MovieID).Return(errors.New("error"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/v1/ratings/1", nil)
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
