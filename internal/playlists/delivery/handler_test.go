package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"errors"

	uuid "github.com/satori/go.uuid"

	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/logger"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/middleware"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	playlistsMock "github.com/go-park-mail-ru/2021_1_kekEnd/internal/playlists/mocks"
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

	playlistsUC := playlistsMock.NewMockUseCase(ctrl)
	usersUC := usersMock.NewMockUseCase(ctrl)
	delivery := mocks.NewMockDelivery(ctrl)

	authMiddleware := middleware.NewAuthMiddleware(usersUC, delivery)
	api := r.Group("/api")
	v1 := api.Group("/v1")
	RegisterHttpEndpoints(v1, playlistsUC, usersUC, authMiddleware, lg)

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

	playlist := models.Playlist{
		ID:       "1",
		Name:     "BestMovies",
		IsShared: false,
	}
	newBody, err := json.Marshal(playlist)
	assert.NoError(t, err)

	playlistMovieData := PlaylistMovie{
		MovieID: "1",
	}
	newBody2, err := json.Marshal(playlistMovieData)
	assert.NoError(t, err)

	playlistUserData := PlaylistUser{
		Username: "let_robots_reign",
	}
	newBody3, err := json.Marshal(playlistUserData)
	assert.NoError(t, err)

	playlistMovieData2 := PlaylistMovie{
		MovieID: "qwe",
	}
	newBody4, err := json.Marshal(playlistMovieData2)
	assert.NoError(t, err)

	usersUC.EXPECT().GetUser(user.Username).Return(user, nil).AnyTimes()
	UUID := uuid.NewV4().String()
	delivery.EXPECT().GetUser(UUID).Return(user.Username, nil).AnyTimes()

	cookie := &http.Cookie{
		Name:  "session_id",
		Value: UUID,
	}

	t.Run("CreatePlaylist", func(t *testing.T) {
		playlistsUC.EXPECT().CreatePlaylist(user.Username, playlist.Name, playlist.IsShared).Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/playlists", bytes.NewBuffer(newBody))
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("GetPlaylist", func(t *testing.T) {
		playlistsUC.EXPECT().GetPlaylist(1).Return(nil, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/playlist/1", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("GetPlaylistsInfo", func(t *testing.T) {
		playlistsUC.EXPECT().GetPlaylistsInfo(user.Username, 1).Return(nil, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/playlists/movies/1", nil)
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("GetPlaylists", func(t *testing.T) {
		playlistsUC.EXPECT().GetPlaylists(user.Username).Return(nil, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/playlists/users/let_robots_reign", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("UpdatePlaylist", func(t *testing.T) {
		playlistsUC.EXPECT().UpdatePlaylist(user.Username, 1, playlist.Name, playlist.IsShared).Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/api/v1/playlists", bytes.NewBuffer(newBody))
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("DeletePlaylist", func(t *testing.T) {
		playlistsUC.EXPECT().DeletePlaylist(user.Username, 1).Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/v1/playlists/1", nil)
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("AddMovieToPlaylist", func(t *testing.T) {
		playlistsUC.EXPECT().AddMovieToPlaylist(user.Username, 1, 1).Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/playlists/1/movie", bytes.NewBuffer(newBody2))
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("DeleteMovieFromPlaylist", func(t *testing.T) {
		playlistsUC.EXPECT().DeleteMovieFromPlaylist(user.Username, 1, 1).Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/v1/playlists/1/movie", bytes.NewBuffer(newBody2))
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("AddUserToPlaylist", func(t *testing.T) {
		playlistsUC.EXPECT().AddUserToPlaylist(user.Username, 1, "let_robots_reign").Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/playlists/1/user", bytes.NewBuffer(newBody3))
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("DeleteUserFromPlaylist", func(t *testing.T) {
		playlistsUC.EXPECT().DeleteUserFromPlaylist(user.Username, 1, "let_robots_reign").Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/v1/playlists/1/user", bytes.NewBuffer(newBody3))
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("CreatePlaylistError", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/playlists", nil)
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("CreatePlaylistError2", func(t *testing.T) {
		playlistsUC.EXPECT().CreatePlaylist(user.Username, playlist.Name, playlist.IsShared).Return(errors.New("error"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/playlists", bytes.NewBuffer(newBody))
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("GetPlaylistError", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/playlist/:playlist_id", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("GetPlaylistError2", func(t *testing.T) {
		playlistsUC.EXPECT().GetPlaylist(1).Return(nil, errors.New("error"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/playlist/1", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("GetPlaylistsInfoError", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/playlists/movies/:movie_id", nil)
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("GetPlaylistsInfoError2", func(t *testing.T) {
		playlistsUC.EXPECT().GetPlaylistsInfo(user.Username, 1).Return(nil, errors.New("error"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/playlists/movies/1", nil)
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("GetPlaylistsError", func(t *testing.T) {
		playlistsUC.EXPECT().GetPlaylists(user.Username).Return(nil, errors.New("error"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/playlists/users/let_robots_reign", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("UpdatePlaylistError", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/api/v1/playlists", nil)
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("UpdatePlaylistError2", func(t *testing.T) {
		playlistsUC.EXPECT().UpdatePlaylist(user.Username, 1, playlist.Name, playlist.IsShared).Return(errors.New("error"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/api/v1/playlists", bytes.NewBuffer(newBody))
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("DeletePlaylistError", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/v1/playlists/:playlist_id", nil)
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("DeletePlaylistError2", func(t *testing.T) {
		playlistsUC.EXPECT().DeletePlaylist(user.Username, 1).Return(errors.New("error"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/v1/playlists/1", nil)
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("AddMovieToPlaylistError", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/playlists/1/movie", nil)
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("AddMovieToPlaylistError2", func(t *testing.T) {
		playlistsUC.EXPECT().AddMovieToPlaylist(user.Username, 1, 1).Return(errors.New("error"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/playlists/1/movie", bytes.NewBuffer(newBody2))
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("AddMovieToPlaylistError3", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/playlists/:playlist_id/movie", bytes.NewBuffer(newBody2))
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("AddMovieToPlaylistError4", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/playlists/1/movie", bytes.NewBuffer(newBody4))
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("DeleteMovieFromPlaylistError", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/v1/playlists/1/movie", nil)
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("DeleteMovieFromPlaylistError2", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/v1/playlists/:playlist_id/movie", bytes.NewBuffer(newBody2))
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("DeleteMovieFromPlaylistError3", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/v1/playlists/1/movie", bytes.NewBuffer(newBody4))
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("DeleteMovieFromPlaylistError4", func(t *testing.T) {
		playlistsUC.EXPECT().DeleteMovieFromPlaylist(user.Username, 1, 1).Return(errors.New("error"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/v1/playlists/1/movie", bytes.NewBuffer(newBody2))
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("AddUserToPlaylistError", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/playlists/1/user", nil)
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("AddUserToPlaylistError2", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/playlists/:playlist_id/user", bytes.NewBuffer(newBody3))
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("AddUserToPlaylistError3", func(t *testing.T) {
		playlistsUC.EXPECT().AddUserToPlaylist(user.Username, 1, "let_robots_reign").Return(errors.New("error"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/playlists/1/user", bytes.NewBuffer(newBody3))
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("DeleteUserFromPlaylistError", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/v1/playlists/1/user", nil)
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("DeleteUserFromPlaylistError2", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/v1/playlists/:playlist_id/user", bytes.NewBuffer(newBody3))
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("DeleteUserFromPlaylistError3", func(t *testing.T) {
		playlistsUC.EXPECT().DeleteUserFromPlaylist(user.Username, 1, "let_robots_reign").Return(errors.New("error"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/v1/playlists/1/user", bytes.NewBuffer(newBody3))
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
