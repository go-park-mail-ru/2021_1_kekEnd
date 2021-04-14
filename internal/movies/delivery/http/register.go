package http

import (
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/movies"
)

func RegisterHttpEndpoints(router *gin.Engine, moviesUC movies.UseCase) {
	handler := NewHandler(moviesUC)

	router.POST("/movies", handler.CreateMovie)
	router.GET("/movies", handler.GetMovies)
	router.GET("/movies/:id", handler.GetMovie)
	router.GET("/genres", handler.GetGenres)
}
