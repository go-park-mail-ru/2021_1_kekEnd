package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/logger"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/middleware"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/proto"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/services/sessions/mocks"
	usersMock "github.com/go-park-mail-ru/2021_1_kekEnd/internal/users/mocks"
	"github.com/golang/mock/gomock"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type MockFileServerClient struct {
}

func (mockFileServer *MockFileServerClient) Upload(ctx context.Context, opts ...grpc.CallOption) (proto.FileServerHandler_UploadClient, error) {
	return nil, nil
}

func NewMockFileServerClient() MockFileServerClient {
	return MockFileServerClient{}
}

func TestHandlers(t *testing.T) {
	r := gin.Default()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	lg := logger.NewAccessLogger()

	usersUC := usersMock.NewMockUseCase(ctrl)
	sessionsUC := mocks.NewMockUseCase(ctrl)
	delivery := mocks.NewMockDelivery(ctrl)

	authMiddleware := middleware.NewAuthMiddleware(usersUC, delivery)

	fileServerMock := NewMockFileServerClient()

	RegisterHttpEndpoints(r, usersUC, delivery, authMiddleware, &fileServerMock, lg)

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
		Avatar:        "http://localhost:8080/avatars/default.jpeg",
		MoviesWatched: new(uint),
		ReviewsNumber: new(uint),
		Subscriptions: new(uint),
		Subscribers:   new(uint),
	}

	UUID := uuid.NewV4().String()
	testErr := errors.New("")

	t.Run("CreateUser", func(t *testing.T) {
		sessionsUC.
			EXPECT().
			Create(user.Username, 240*time.Hour).
			Return(UUID, nil).AnyTimes()

		delivery.EXPECT().Create(user.Username, 240*time.Hour).Return(UUID, nil).AnyTimes()

		sessionID, err := delivery.Create(user.Username, 240*time.Hour)
		assert.NoError(t, err)
		assert.Equal(t, UUID, sessionID)

		usersUC.EXPECT().CreateUser(user).Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("TestLogout", func(t *testing.T) {
		mockUser := &models.User{
			Username:      "let_robots_reign",
			Email:         "sample@ya.ru",
			Password:      "1234",
			Avatar:        "http://localhost:8080/avatars/default.jpeg",
			MoviesWatched: new(uint),
			ReviewsNumber: new(uint),
			Subscribers:   new(uint),
			Subscriptions: new(uint),
		}

		cookie := &http.Cookie{
			Name:  "session_id",
			Value: UUID,
		}

		usersUC.
			EXPECT().
			GetUser(user.Username).
			Return(mockUser, nil).AnyTimes()

		sessionsUC.
			EXPECT().
			GetUser(UUID).
			Return(user.Username, nil).AnyTimes()

		sessionsUC.
			EXPECT().
			Delete(cookie.Value).
			Return(nil).AnyTimes()

		delivery.EXPECT().GetUser(UUID).Return(user.Username, nil).AnyTimes()
		delivery.EXPECT().Delete(cookie.Value).Return(nil).AnyTimes()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/sessions", nil)
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("GetUser", func(t *testing.T) {
		mockUser := &models.User{
			Username:      "let_robots_reign",
			Email:         "sample@ya.ru",
			Password:      "1234",
			Avatar:        "http://localhost:8080/avatars/default.jpeg",
			MoviesWatched: new(uint),
			ReviewsNumber: new(uint),
			Subscribers:   new(uint),
			Subscriptions: new(uint),
		}

		cookie := &http.Cookie{
			Name:  "session_id",
			Value: UUID,
		}

		usersUC.EXPECT().GetUser(user.Username).Return(mockUser, nil).AnyTimes()

		sessionsUC.
			EXPECT().
			GetUser(UUID).
			Return(user.Username, nil).AnyTimes()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/users", nil)
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("TestUpdateUser", func(t *testing.T) {
		newMockUser := models.User{
			Username:      "let_robots_reign",
			Email:         "corrected@ya.ru",
			Password:      "1234",
			Avatar:        "http://localhost:8080/avatars/default.jpeg",
			MoviesWatched: new(uint),
			ReviewsNumber: new(uint),
			Subscribers:   new(uint),
			Subscriptions: new(uint),
		}

		body, err := json.Marshal(newMockUser)
		assert.NoError(t, err)

		cookie := &http.Cookie{
			Name:  "session_id",
			Value: UUID,
		}

		usersUC.EXPECT().UpdateUser(user, newMockUser).Return(&newMockUser, nil).AnyTimes()

		delivery.EXPECT().GetUser(UUID).Return(user.Username, nil).AnyTimes()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/users", bytes.NewBuffer(body))
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("TestSubscribe", func(t *testing.T) {
		cookie := &http.Cookie{
			Name:  "session_id",
			Value: UUID,
		}

		usersUC.EXPECT().Subscribe(user.Username, user.Username).Return(nil)

		sessionsUC.
			EXPECT().
			GetUser(UUID).
			Return(user.Username, nil).AnyTimes()

		delivery.EXPECT().GetUser(UUID).Return(user.Username, nil).AnyTimes()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/subscriptions/let_robots_reign", bytes.NewBuffer(body))
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("TestSubscribe-FAIL", func(t *testing.T) {
		cookie := &http.Cookie{
			Name:  "session_id",
			Value: UUID,
		}

		usersUC.EXPECT().Subscribe(user.Username, user.Username).Return(testErr)

		sessionsUC.
			EXPECT().
			GetUser(UUID).
			Return(user.Username, nil).AnyTimes()

		delivery.EXPECT().GetUser(UUID).Return(user.Username, nil).AnyTimes()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/subscriptions/let_robots_reign", bytes.NewBuffer(body))
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("TestUnsubscribe", func(t *testing.T) {
		cookie := &http.Cookie{
			Name:  "session_id",
			Value: UUID,
		}

		usersUC.EXPECT().Unsubscribe(user.Username, user.Username).Return(nil)

		sessionsUC.
			EXPECT().
			GetUser(UUID).
			Return(user.Username, nil).AnyTimes()

		delivery.EXPECT().GetUser(UUID).Return(user.Username, nil).AnyTimes()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/subscriptions/let_robots_reign", bytes.NewBuffer(body))
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("TestUnsubscribe-FAIL", func(t *testing.T) {
		cookie := &http.Cookie{
			Name:  "session_id",
			Value: UUID,
		}

		usersUC.EXPECT().Unsubscribe(user.Username, user.Username).Return(testErr)

		sessionsUC.
			EXPECT().
			GetUser(UUID).
			Return(user.Username, nil).AnyTimes()

		delivery.EXPECT().GetUser(UUID).Return(user.Username, nil).AnyTimes()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/subscriptions/let_robots_reign", bytes.NewBuffer(body))
		req.AddCookie(cookie)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("TestGetSubscribers", func(t *testing.T) {
		page := 1
		subs := make([]models.UserNoPassword, 0)
		usersUC.EXPECT().GetSubscribers(page, user.Username).Return(1, subs, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/subscribers/let_robots_reign", bytes.NewBuffer(body))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("TestGetSubscribers-FAIL", func(t *testing.T) {
		page := 1
		usersUC.EXPECT().GetSubscribers(page, user.Username).Return(1, nil, testErr)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/subscribers/let_robots_reign", bytes.NewBuffer(body))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("TestGetSubscriptions", func(t *testing.T) {
		page := 1
		subs := make([]models.UserNoPassword, 0)
		usersUC.EXPECT().GetSubscriptions(page, user.Username).Return(1, subs, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/subscriptions/let_robots_reign", bytes.NewBuffer(body))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("TestGetSubscriptions-FAIL", func(t *testing.T) {
		page := 1
		usersUC.EXPECT().GetSubscriptions(page, user.Username).Return(1, nil, testErr)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/subscriptions/let_robots_reign", bytes.NewBuffer(body))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

}
