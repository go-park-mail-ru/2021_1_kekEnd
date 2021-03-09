package http

import (
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/sessions"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/users"
)

func RegisterHttpEndpoints(router *gin.Engine, usersUC users.UseCase, sessions sessions.Delivery) {
	handler := NewHandler(usersUC, sessions)

	router.POST("/users", handler.CreateUser)
	router.POST("/users/login", handler.Login)
	router.GET("/users/:username", handler.GetUser)
	router.PUT("/users/:username", handler.UpdateUser)

	router.GET("/checkAuth", handler.CheckAuth)
	router.POST("/logout", handler.Logout)
}
