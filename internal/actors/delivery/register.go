package delivery

import (
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/actors"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/middleware"
)

func RegisterHttpEndpoints(router *gin.Engine, actorsUC actors.UseCase, auth middleware.Auth) {
	handler := NewHandler(actorsUC)

	router.POST("/actors", auth.CheckAuth(), handler.CreateActor)
	router.GET("/actors/:actor_id", handler.GetActor)
	router.PUT("/actors/:actor_id", auth.CheckAuth(), handler.EditActor)
}
