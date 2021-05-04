package http

import (
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/logger"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/middleware"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/sessions"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/users"
)

func RegisterHttpEndpoints(router *gin.Engine, usersUC users.UseCase, sessions sessions.Delivery,
	authMiddleware middleware.Auth, Log *logger.Logger) {
	handler := NewHandler(usersUC, sessions, Log)

	router.POST("/users", handler.CreateUser)
	router.POST("/users/avatar", authMiddleware.RequireAuth(), handler.UploadAvatar)
	router.GET("/users", authMiddleware.RequireAuth(), handler.GetUser)
	router.PUT("/users", authMiddleware.RequireAuth(), handler.UpdateUser)
	router.DELETE("/sessions", authMiddleware.RequireAuth(), handler.Logout)
	router.POST("/sessions", handler.Login)

	router.GET("/subscriptions/:user_id", handler.GetSubscriptions)
	router.POST("/subscriptions/:user_id", authMiddleware.RequireAuth(), handler.Subscribe)
	router.DELETE("/subscriptions/:user_id", authMiddleware.RequireAuth(), handler.Unsubscribe)
	router.GET("/subscribers/:user_id", handler.GetSubscribers)
	router.GET("/feed", authMiddleware.CheckAuth(), handler.GetFeed)
}
