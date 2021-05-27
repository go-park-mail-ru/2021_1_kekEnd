package http

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
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/movies/mocks"
	sessionMock "github.com/go-park-mail-ru/2021_1_kekEnd/internal/services/sessions/mocks"
	usersMock "github.com/go-park-mail-ru/2021_1_kekEnd/internal/users/mocks"
	"github.com/golang/mock/gomock"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestMoviesHandlers(t *testing.T) {
	r := gin.Default()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	moviesUC := mocks.NewMockUseCase(ctrl)
	usersUC := usersMock.NewMockUseCase(ctrl)
	delivery := sessionMock.NewMockDelivery(ctrl)
	authMiddleware := middleware.NewAuthMiddleware(usersUC, delivery)
	lg := logger.NewAccessLogger()
	api := r.Group("/api")
	v1 := api.Group("/v1")
	RegisterHTTPEndpoints(v1, moviesUC, authMiddleware, lg)

	user := &models.User{
		Username:      "let_robots_reign",
		Email:         "sample@mail.ru",
		Password:      "123",
		Avatar:        "http://localhost:8080/avatars/default.jpeg",
		MoviesWatched: new(uint),
		ReviewsNumber: new(uint),
	}

	movie := &models.Movie{
		ID:          "7",
		Title:       "Yet another movie",
		Description: "Generic description",
		Rating:      8.5,
		RatingCount: 1,
	}

	usersUC.EXPECT().GetUser(user.Username).Return(user, nil).AnyTimes()
	UUID := uuid.NewV4().String()
	delivery.EXPECT().GetUser(UUID).Return(user.Username, nil).AnyTimes()

	body, err := json.Marshal(movie)
	assert.NoError(t, err)

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: UUID,
	}

	t.Run("CreateMovie", func(t *testing.T) {
		moviesUC.EXPECT().CreateMovie(movie).Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/movies", bytes.NewBuffer(body))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("GetMovie", func(t *testing.T) {
		moviesUC.EXPECT().GetMovie(movie.ID, "").Return(movie, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/movies/"+movie.ID, bytes.NewBuffer(body))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("GetBestMovies", func(t *testing.T) {
		moviesUC.EXPECT().GetBestMovies(1, "").Return(1, []*models.Movie{
			movie,
		}, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/movies?category=best&page=1", bytes.NewBuffer(body))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("GetGenres", func(t *testing.T) {
		moviesUC.EXPECT().GetAllGenres().Return([]string{"драма"}, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/genres", bytes.NewBuffer(body))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("GetMoviesByGenres", func(t *testing.T) {
		moviesUC.EXPECT().GetMoviesByGenres([]string{"драма"}, 1, "").Return(1, []*models.Movie{
			movie,
		}, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/movies?category=genre&filter=драма&page=1", bytes.NewBuffer(body))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("MarkWatched", func(t *testing.T) {
		moviesUC.EXPECT().MarkWatched(*user, 1).Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/movies/1/watch", bytes.NewBuffer(body))
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("MarkUnwatched", func(t *testing.T) {
		moviesUC.EXPECT().MarkUnwatched(*user, 1).Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/v1/movies/1/watch", bytes.NewBuffer(body))
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("CreateMovieError", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/movies", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("CreateMovieError2", func(t *testing.T) {
		moviesUC.EXPECT().CreateMovie(movie).Return(errors.New("error"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/movies", bytes.NewBuffer(body))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("GetMovieError", func(t *testing.T) {
		moviesUC.EXPECT().GetMovie(movie.ID, "").Return(movie, errors.New("error"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/movies/"+movie.ID, nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("GetBestMoviesError", func(t *testing.T) {
		moviesUC.EXPECT().GetBestMovies(1, "").Return(1, []*models.Movie{
			movie,
		}, errors.New("error"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/movies?category=best&page=1", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("GetGenresError", func(t *testing.T) {
		moviesUC.EXPECT().GetAllGenres().Return(nil, errors.New("error"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/genres", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("GetMoviesByGenresError", func(t *testing.T) {
		moviesUC.EXPECT().GetMoviesByGenres([]string{"драма"}, 1, "").Return(1, nil, errors.New("error"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/movies?category=genre&filter=драма&page=1", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("MarkWatchedError", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/movies/:movie_id/watch", nil)
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("MarkWatchedError2", func(t *testing.T) {
		moviesUC.EXPECT().MarkWatched(*user, 1).Return(errors.New("error"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/movies/1/watch", nil)
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("MarkUnwatchedError", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/v1/movies/:movie_id/watch", nil)
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("MarkUnwatchedError2", func(t *testing.T) {
		moviesUC.EXPECT().MarkUnwatched(*user, 1).Return(errors.New("error"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/v1/movies/1/watch", nil)
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
