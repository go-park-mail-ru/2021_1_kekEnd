package actors

import (
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/actors"
	"net/http"
)

type Handler struct {
	useCase actors.UseCase
}

func NewHandler(useCase actors.UseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

func (h *Handler) GetActor(ctx *gin.Context) {
	id := ctx.Param("id")

	actor, err := h.useCase.GetActor(id)

	if err != nil {
		ctx.AbortWithStatus(http.StatusNotFound) // 404
		return
	}

	ctx.JSON(http.StatusOK, actor)
}
