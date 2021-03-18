package http

import (
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/reviews"
	_const "github.com/go-park-mail-ru/2021_1_kekEnd/pkg/const"
	"net/http"
)

type Handler struct {
	useCase reviews.UseCase
}

func NewHandler(useCase reviews.UseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

func (h *Handler) CreateReview(ctx *gin.Context) {
	review := new(models.Review)
	err := ctx.BindJSON(review)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	user, ok := ctx.Get(_const.UserKey)
	if !ok {
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	userModel, ok := user.(models.User)
	if !ok {
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	review.Author = userModel.Username

	err = h.useCase.CreateReview(review)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	ctx.Status(http.StatusCreated) // 201
}
