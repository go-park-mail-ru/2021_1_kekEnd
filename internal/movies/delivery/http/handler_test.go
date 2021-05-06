package http

import (
	"bytes"
	"encoding/json"
	"errors"
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
	"net/http"
	"net/http/httptest"
	"testing"
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
	RegisterHttpEndpoints(r, moviesUC, authMiddleware, lg)

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

	t.Run("CreateMovie", func(t *testing.T) {
		moviesUC.EXPECT().CreateMovie(movie).Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/movies", bytes.NewBuffer(body))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("GetMovie", func(t *testing.T) {
		moviesUC.EXPECT().GetMovie(movie.ID, "").Return(movie, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/movies/"+movie.ID, bytes.NewBuffer(body))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("GetMovieError", func(t *testing.T) {
		wrongID := "42"
		moviesUC.EXPECT().GetMovie(wrongID, "").Return(nil, errors.New("movie not found"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/movies/"+wrongID, bytes.NewBuffer(body))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("GetBestMovies", func(t *testing.T) {
		moviesUC.EXPECT().GetBestMovies(1, "").Return(1, []*models.Movie{
			movie,
		}, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/movies?category=best&page=1", bytes.NewBuffer(body))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("GetGenres", func(t *testing.T) {
		moviesUC.EXPECT().GetAllGenres().Return([]string{"драма"}, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/genres", bytes.NewBuffer(body))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("GetMoviesByGenres", func(t *testing.T) {
		moviesUC.EXPECT().GetMoviesByGenres([]string{"драма"}, 1, "").Return(1, []*models.Movie{
			movie,
		}, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/movies?category=genre&filter=драма&page=1", bytes.NewBuffer(body))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
