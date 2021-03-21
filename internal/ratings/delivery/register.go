package ratings

import (
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/ratings"
	)

func RegisterHttpEndpoints(router *gin.Engine, ratingsUC ratings.UseCase) {

	h := NewHandler(ratingsUC)

	//router.POST("/movies", handler)
	//router.GET("/movies/:id", handler.GetMovie)
}