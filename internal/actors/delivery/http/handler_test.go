package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	uuid "github.com/satori/go.uuid"

	"github.com/gin-gonic/gin"
	actorsMock "github.com/go-park-mail-ru/2021_1_kekEnd/internal/actors/mocks"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/logger"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/middleware"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
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

	actorsUC := actorsMock.NewMockUseCase(ctrl)
	usersUC := usersMock.NewMockUseCase(ctrl)
	delivery := mocks.NewMockDelivery(ctrl)

	authMiddleware := middleware.NewAuthMiddleware(usersUC, delivery)

	RegisterHttpEndpoints(r, actorsUC, authMiddleware, lg)

	user := &models.User{
		Username:      "let_robots_reign",
		Email:         "sample@mail.ru",
		Password:      "123",
		Avatar:        "http://localhost:8080/avatars/default.jpeg",
		MoviesWatched: new(uint),
		ReviewsNumber: new(uint),
	}

	actor := models.Actor{
		ID:   "1",
		Name: "Tom Cruise",
	}

	usersUC.EXPECT().GetUser(user.Username).Return(user, nil).AnyTimes()
	UUID := uuid.NewV4().String()
	delivery.EXPECT().GetUser(UUID).Return(user.Username, nil).AnyTimes()

	t.Run("GetActor", func(t *testing.T) {
		actorsUC.EXPECT().GetActor(actor.ID, "").Return(actor, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/actors/1", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
