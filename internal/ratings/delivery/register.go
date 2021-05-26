package ratings

import (
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/logger"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/middleware"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/ratings"
)

func RegisterHttpEndpoints(router *gin.RouterGroup, ratingsUC ratings.UseCase, authMiddleware middleware.Auth,
	Log *logger.Logger) {
	handler := NewHandler(ratingsUC, Log)

	router.POST("/ratings", authMiddleware.CheckAuth(true), handler.CreateRating)
	router.GET("/ratings/:movie_id", authMiddleware.CheckAuth(true), handler.GetRating)
	router.PUT("/ratings", authMiddleware.CheckAuth(true), handler.UpdateRating)
	router.DELETE("/ratings/:movie_id", authMiddleware.CheckAuth(true), handler.DeleteRating)
}
