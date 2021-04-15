package http
//
//import (
//	"bytes"
//	"encoding/json"
//	"github.com/gin-gonic/gin"
//	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/middleware"
//	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
//	sessionsMock "github.com/go-park-mail-ru/2021_1_kekEnd/internal/sessions"
//	sessions "github.com/go-park-mail-ru/2021_1_kekEnd/internal/sessions/delivery"
//	"github.com/golang/mock/gomock"
//	uuid "github.com/satori/go.uuid"
//	"github.com/stretchr/testify/assert"
//	"net/http"
//	"net/http/httptest"
//	"testing"
//	"time"
//)
//
//func TestHandlers(t *testing.T) {
//	r := gin.Default()
//
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	usersUC := usersMock.NewMockUseCase(ctrl)
//	sessionsUC := sessionsMock.NewMockUseCase(ctrl)
//	delivery := sessions.NewDelivery(sessionsUC)
//
//	authMiddleware := middleware.NewAuthMiddleware(usersUC, delivery)
//
//	RegisterHttpEndpoints(r, usersUC, delivery, authMiddleware)
//
//	createBody := &signupData{
//		Username: "let_robots_reign",
//		Email:    "sample@ya.ru",
//		Password: "1234",
//	}
//
//	body, err := json.Marshal(createBody)
//	assert.NoError(t, err)
//
//	user := &models.User{
//		Username:      createBody.Username,
//		Email:         createBody.Email,
//		Password:      createBody.Password,
//		Avatar:        "http://localhost:8080/avatars/default.jpeg",
//		MoviesWatched: new(uint),
//		ReviewsNumber: new(uint),
//	}
//
//	UUID := uuid.NewV4().String()
//
//	t.Run("CreateUser", func(t *testing.T) {
//		sessionsUC.
//			EXPECT().
//			Create(user.Username, 240*time.Hour).
//			Return(UUID, nil).AnyTimes()
//
//		sessionID, err := delivery.Create(user.Username, 240*time.Hour)
//		assert.NoError(t, err)
//		assert.Equal(t, UUID, sessionID)
//
//		usersUC.EXPECT().CreateUser(user).Return(nil)
//
//		w := httptest.NewRecorder()
//		req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
//		r.ServeHTTP(w, req)
//
//		assert.Equal(t, http.StatusCreated, w.Code)
//	})
//
//	t.Run("GetUser", func(t *testing.T) {
//		mockUser := &models.User{
//			Username:      "let_robots_reign",
//			Email:         "sample@ya.ru",
//			Password:      "1234",
//			Avatar:        "http://localhost:8080/avatars/default.jpeg",
//			MoviesWatched: new(uint),
//			ReviewsNumber: new(uint),
//		}
//
//		cookie := &http.Cookie{
//			Name:  "session_id",
//			Value: UUID,
//		}
//
//		usersUC.EXPECT().GetUser(user.Username).Return(mockUser, nil).AnyTimes()
//
//		sessionsUC.
//			EXPECT().
//			Check(UUID).
//			Return(user.Username, nil).AnyTimes()
//
//		w := httptest.NewRecorder()
//		req, _ := http.NewRequest("GET", "/users", nil)
//		req.AddCookie(cookie)
//		r.ServeHTTP(w, req)
//
//		assert.Equal(t, http.StatusOK, w.Code)
//	})
//
//	t.Run("TestUpdateUser", func(t *testing.T) {
//		newMockUser := models.User{
//			Username:      "let_robots_reign",
//			Email:         "corrected@ya.ru",
//			Password:      "1234",
//			Avatar:        "http://localhost:8080/avatars/default.jpeg",
//			MoviesWatched: new(uint),
//			ReviewsNumber: new(uint),
//		}
//
//		body, err := json.Marshal(newMockUser)
//		assert.NoError(t, err)
//
//		cookie := &http.Cookie{
//			Name:  "session_id",
//			Value: UUID,
//		}
//
//		usersUC.EXPECT().UpdateUser(user, newMockUser).Return(&newMockUser, nil)
//
//		sessionsUC.
//			EXPECT().
//			Check(UUID).
//			Return(user.Username, nil).AnyTimes()
//
//		w := httptest.NewRecorder()
//		req, _ := http.NewRequest("PUT", "/users", bytes.NewBuffer(body))
//		req.AddCookie(cookie)
//		r.ServeHTTP(w, req)
//
//		assert.Equal(t, http.StatusOK, w.Code)
//	})
//}
