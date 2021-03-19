package http

import (
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/middleware"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/reviews"
)

func RegisterHttpEndpoints(router *gin.Engine, reviewsUC reviews.UseCase, authMiddleware middleware.Auth) {
	handler := NewHandler(reviewsUC)

	router.POST("/users/reviews", authMiddleware.CheckAuth(), handler.CreateReview)
	router.GET("/movies/:id/reviews", handler.GetMovieReviews)
	router.GET("/users/reviews", authMiddleware.CheckAuth(), handler.GetUserReviews)
	router.GET("/users/movies/:id/reviews", authMiddleware.CheckAuth(), handler.GetUserReviewForMovie)
}
