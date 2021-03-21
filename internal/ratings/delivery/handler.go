package ratings

import (
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/ratings"
	_const "github.com/go-park-mail-ru/2021_1_kekEnd/pkg/const"
	"net/http"
)

type Handler struct {
	usecase ratings.UseCase
}

func NewHandler(usecase ratings.UseCase) *Handler {
	return &Handler{
		usecase: usecase,
	}
}

func (h *Handler) CreateRating(ctx *gin.Context) {
	var score uint = 0
	err := ctx.BindJSON(score)

	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
	}

	user, ok := ctx.Get(_const.UserKey)
	if !ok {
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
	}

	username := user.(models.User).Username
	movieID := ctx.Param("movie_id")
	err = h.usecase.CreateRating(username, movieID, score)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
	}

	ctx.Status(http.StatusOK)
}

func (h *Handler) GetRating(ctx *gin.Context) {
	user, ok := ctx.Get(_const.UserKey)
	if !ok {
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
	}

	username := user.(models.User).Username
	movieID := ctx.Param("movie_id")
	rating, err := h.usecase.GetRating(username, movieID)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
	}

	ctx.JSON(http.StatusOK, rating)
}

func (h *Handler) DeleteRating(ctx *gin.Context) {
	user, ok := ctx.Get(_const.UserKey)
	if !ok {
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
	}

	username := user.(models.User).Username
	movieID := ctx.Param("movie_id")
	err := h.usecase.DeleteRating(username, movieID)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
	}

	ctx.Status(http.StatusOK)
}
