package http

import (
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2021_1_kekEnd/auth"
)

func RegisterHttpEndpoints(router *gin.Engine, authUC auth.UseCase) {
	handler := NewHandler(authUC)

	authEndpoints := router.Group("/auth")
	{
		authEndpoints.POST("/signup", handler.SignUp)
		authEndpoints.POST("/login", handler.Login)
	}
}
