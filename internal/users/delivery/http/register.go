package http

import (
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/users"
)

func RegisterHttpEndpoints(router *gin.Engine, usersUC users.UseCase) {
	handler := NewHandler(usersUC)

	router.POST("/users", handler.CreateUser)
	router.GET("/users/:username", handler.GetUser)
	router.PUT("/users/:username", handler.UpdateUser)
}
