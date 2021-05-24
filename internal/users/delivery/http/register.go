package http

import (
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/logger"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/middleware"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/proto"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/services/sessions"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/users"
)

func RegisterHttpEndpoints(router *gin.RouterGroup, usersUC users.UseCase, sessions sessions.Delivery,
	authMiddleware middleware.Auth, fileServer proto.FileServerHandlerClient, Log *logger.Logger) {
	handler := NewHandler(usersUC, sessions, fileServer, Log)

	router.POST("/users", handler.CreateUser)
	router.POST("/users/avatar", authMiddleware.CheckAuth(true), handler.UploadAvatar)
	router.GET("/users", authMiddleware.CheckAuth(true), handler.GetCurrentUser)
	router.GET("/user/:username", handler.GetUser)
	router.PUT("/users", authMiddleware.CheckAuth(true), handler.UpdateUser)
	router.DELETE("/sessions", authMiddleware.CheckAuth(true), handler.Logout)
	router.POST("/sessions", handler.Login)

	router.GET("/subscriptions/:username", handler.GetSubscriptions)
	router.POST("/subscriptions/:username", authMiddleware.CheckAuth(true), handler.Subscribe)
	router.DELETE("/subscriptions/:username", authMiddleware.CheckAuth(true), handler.Unsubscribe)
	router.GET("/subscriptions/:username/check", authMiddleware.CheckAuth(true), handler.IsSubscribed)
	router.GET("/subscribers/:username", handler.GetSubscribers)
	router.GET("/feed", authMiddleware.CheckAuth(true), handler.GetFeed)
}
