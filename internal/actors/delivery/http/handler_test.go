package http

import (
	"github.com/gin-gonic/gin"
	actorsMock "github.com/go-park-mail-ru/2021_1_kekEnd/internal/actors/mocks"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/logger"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/middleware"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	sessionsMock "github.com/go-park-mail-ru/2021_1_kekEnd/internal/sessions"
	sessions "github.com/go-park-mail-ru/2021_1_kekEnd/internal/sessions/delivery"
	usersMock "github.com/go-park-mail-ru/2021_1_kekEnd/internal/users/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlers(t *testing.T) {
	r := gin.Default()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	lg := logger.NewAccessLogger()

	actorsUC := actorsMock.NewMockUseCase(ctrl)
	usersUC := usersMock.NewMockUseCase(ctrl)
	sessionsUC := sessionsMock.NewMockUseCase(ctrl)
	delivery := sessions.NewDelivery(sessionsUC, lg)

	authMiddleware := middleware.NewAuthMiddleware(usersUC, delivery)

	RegisterHttpEndpoints(r, actorsUC, authMiddleware, lg)

	actor := models.Actor{
		ID:   "1",
		Name: "Tom Cruise",
	}

	t.Run("GetActor", func(t *testing.T) {
		actorsUC.EXPECT().GetActor(actor.ID).Return(actor, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/actors/1", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}