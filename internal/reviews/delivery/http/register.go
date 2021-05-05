package http

import (
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/logger"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/middleware"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/reviews"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/users"
)

func RegisterHttpEndpoints(router *gin.Engine, reviewsUC reviews.UseCase, usersUC users.UseCase,
	authMiddleware middleware.Auth, Log *logger.Logger) {
	handler := NewHandler(reviewsUC, usersUC, Log)

	router.POST("/users/reviews", authMiddleware.CheckAuth(true), handler.CreateReview)
	router.GET("/movies/:id/reviews", handler.GetMovieReviews)
	router.GET("/user/:user_id/reviews", handler.GetUserReviews)
	router.GET("/users/movies/:id/reviews", authMiddleware.CheckAuth(true), handler.GetUserReviewForMovie)
	router.PUT("/users/movies/:id/reviews", authMiddleware.CheckAuth(true), handler.EditUserReviewForMovie)
	router.DELETE("/users/movies/:id/reviews", authMiddleware.CheckAuth(true), handler.DeleteUserReviewForMovie)
}
