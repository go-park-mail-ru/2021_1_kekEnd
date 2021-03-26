package ratings

import (
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/middleware"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/ratings"
)

func RegisterHttpEndpoints(router *gin.Engine, ratingsUC ratings.UseCase, authMiddleware middleware.Auth) {
	handler := NewHandler(ratingsUC)

	router.POST("/ratings/:movie_id", authMiddleware.CheckAuth(), handler.CreateRating)
	router.GET("/ratings/:movie_id", authMiddleware.CheckAuth(), handler.GetRating)
	router.PUT("/ratings/:movie_id", authMiddleware.CheckAuth(), handler.UpdateRating)
	router.DELETE("/ratings/:movie_id", authMiddleware.CheckAuth(), handler.DeleteRating)
}
