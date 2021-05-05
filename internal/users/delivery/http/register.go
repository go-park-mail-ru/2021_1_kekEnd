package http

import (
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/logger"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/middleware"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/sessions/delivery"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/users"
)

func RegisterHttpEndpoints(router *gin.Engine, usersUC users.UseCase, sessions *delivery.AuthClient,
	authMiddleware middleware.Auth, Log *logger.Logger) {
	handler := NewHandler(usersUC, sessions, Log)

	router.POST("/users", handler.CreateUser)
	router.POST("/users/avatar", authMiddleware.CheckAuth(true), handler.UploadAvatar)
	router.GET("/users", authMiddleware.CheckAuth(true), handler.GetCurrentUser)
	router.GET("/user/:user_id", handler.GetUser)
	router.PUT("/users", authMiddleware.CheckAuth(true), handler.UpdateUser)
	router.DELETE("/sessions", authMiddleware.CheckAuth(true), handler.Logout)
	router.POST("/sessions", handler.Login)

	router.GET("/subscriptions/:user_id", handler.GetSubscriptions)
	router.POST("/subscriptions/:user_id", authMiddleware.CheckAuth(true), handler.Subscribe)
	router.DELETE("/subscriptions/:user_id", authMiddleware.CheckAuth(true), handler.Unsubscribe)
	router.GET("/subscribers/:user_id", handler.GetSubscribers)
	router.GET("/feed", authMiddleware.CheckAuth(true), handler.GetFeed)
}
