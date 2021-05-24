package http

import (
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/logger"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/middleware"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/movies"
)

func RegisterHttpEndpoints(router *gin.RouterGroup, moviesUC movies.UseCase, auth middleware.Auth, Log *logger.Logger) {
	handler := NewHandler(moviesUC, Log)

	router.POST("/movies", handler.CreateMovie)
	router.GET("/movies", auth.CheckAuth(false), handler.GetMovies)
	router.GET("/movies/:id", auth.CheckAuth(false), handler.GetMovie)
	router.POST("/movies/:id/watch", auth.CheckAuth(true), handler.MarkWatched)
	router.DELETE("/movies/:id/watch", auth.CheckAuth(true), handler.MarkUnwatched)
	router.GET("/genres", handler.GetGenres)
}
