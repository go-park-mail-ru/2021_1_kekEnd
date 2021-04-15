package actors

import (
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/actors"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/middleware"
)

func RegisterHttpEndpoints(router *gin.Engine, actorsUC actors.UseCase, authMiddleware middleware.Auth) {
	handler := NewHandler(actorsUC)

	router.GET("/actors/:id", handler.GetActor)
}
