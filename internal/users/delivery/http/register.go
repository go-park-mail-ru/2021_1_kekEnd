package http

import (
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/sessions"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/users"
)

func RegisterHttpEndpoints(router *gin.Engine, usersUC users.UseCase, sessions sessions.Delivery) {
	handler := NewHandler(usersUC, sessions)

	router.POST("/users", handler.CreateUser)
	router.POST("/logout", handler.Logout)
	router.GET("/users/:username", handler.GetUser)
	router.PUT("/users/:username", handler.UpdateUser)
}
