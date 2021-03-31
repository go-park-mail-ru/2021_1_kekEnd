package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/movies/usecase"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMoviesHandlers(t *testing.T) {
	r := gin.Default()
	moviesUC := &usecase.MoviesUseCaseMock{}

	RegisterHttpEndpoints(r, moviesUC)

	movie := &models.Movie{
		ID:          "7",
		Title:       "Yet another movie",
		Description: "Generic description",
	}

	body, err := json.Marshal(movie)
	assert.NoError(t, err)

	t.Run("CreateMovie", func(t *testing.T) {
		moviesUC.On("CreateMovie", movie).Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/movies", bytes.NewBuffer(body))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("GetMovie", func(t *testing.T) {
		moviesUC.On("GetMovie", movie.ID).Return(movie, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/movies/"+movie.ID, bytes.NewBuffer(body))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("GetMovieError", func(t *testing.T) {
		wrongID := "42"
		moviesUC.On("GetMovie", wrongID).Return(nil, errors.New("movie not found"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/movies/"+wrongID, bytes.NewBuffer(body))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
