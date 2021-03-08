package http

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/users/usecase"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateUser(t *testing.T) {
	r := gin.Default()
	uc := &usecase.UsersUseCaseMock{}

	RegisterHttpEndpoints(r, uc)

	createBody := &signupData{
		Username: "let_robots_reign",
		Email:    "sample@ya.ru",
		Password: "1234",
	}

	body, err := json.Marshal(createBody)
	assert.NoError(t, err)

	user := &models.User{
		Username:      createBody.Username,
		Email:         createBody.Email,
		Password:      createBody.Password,
		MoviesWatched: 0,
		ReviewsNumber: 0,
	}

	uc.On("CreateUser", user).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestSuccessfulGetUser(t *testing.T) {
	r := gin.Default()
	uc := &usecase.UsersUseCaseMock{}

	RegisterHttpEndpoints(r, uc)

	username := "let_robots_reign"
	mockUser := &models.User{
		Username:      "let_robots_reign",
		Email:         "sample@ya.ru",
		Password:      "1234",
		MoviesWatched: 0,
		ReviewsNumber: 0,
	}

	uc.On("GetUser", username).Return(mockUser, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/let_robots_reign", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

//func TestUnsuccessfulGetUser(t *testing.T) {
//	r := gin.Default()
//	uc := &usecase.UsersUseCaseMock{}
//
//	RegisterHttpEndpoints(r, uc)
//
//	username := "let_robots_reign"
//
//	uc.On("GetUser", username).Return(nil, errors.New("user not found"))
//
//	w := httptest.NewRecorder()
//	req, _ := http.NewRequest("GET", "/users/let_robots_reign", nil)
//	r.ServeHTTP(w, req)
//
//	assert.Equal(t, http.StatusNotFound, w.Code)
//}

func TestUpdateUser(t *testing.T) {
	r := gin.Default()
	uc := &usecase.UsersUseCaseMock{}

	RegisterHttpEndpoints(r, uc)

	username := "let_robots_reign"
	newMockUser := &models.User{
		Username:      "let_robots_reign",
		Email:         "corrected@ya.ru",
		Password:      "1234",
		MoviesWatched: 0,
		ReviewsNumber: 0,
	}

	body, err := json.Marshal(newMockUser)
	assert.NoError(t, err)

	uc.On("UpdateUser", username, newMockUser).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/users/let_robots_reign", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
