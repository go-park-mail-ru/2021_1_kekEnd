package ratings

import (
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/ratings"
	"net/http"
	"strconv"
)

type Handler struct {
	usecase ratings.UseCase
}

func (h *Handler) NewHandler(usecase ratings.UseCase) *Handler {
	return &Handler{
		usecase: usecase,
	}
}

func (h *Handler) CreateRating(ctx *gin.Context) {
	score, err := strconv.ParseUint(ctx.Param("score"), 10, 64)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
	}

	err = h.usecase.CreateRating(ctx.Param("user_id"), ctx.Param("movie_id"), uint(score))
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
	}

	ctx.Status(http.StatusOK)
}

func (h *Handler) GetRating(ctx *gin.Context) {
	rating, err := h.usecase.GetRating(ctx.Param("user_id"), ctx.Param("movie_id"))
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
	}

	ctx.JSON(http.StatusOK, rating)
}

func (h *Handler) DeleteRating(ctx *gin.Context) {
	err := h.usecase.DeleteRating(ctx.Param("user_id"), ctx.Param("movie_id"))
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
	}

	ctx.Status(http.StatusOK)
}
