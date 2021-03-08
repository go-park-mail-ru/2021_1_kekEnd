package http

import (
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/movies"
	"net/http"
)

type Handler struct {
	useCase movies.UseCase
}

func NewHandler(useCase movies.UseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

func (h *Handler) CreateMovie(ctx *gin.Context) {
	movieData := new(models.Movie)
	err := ctx.BindJSON(movieData)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
	}

	err = h.useCase.CreateMovie(movieData)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
	}

	ctx.Status(http.StatusCreated)
}

func (h *Handler) GetMovie(ctx *gin.Context) {
	movie, err := h.useCase.GetMovie(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatus(http.StatusNotFound) // 404
	}

	ctx.JSON(http.StatusOK, movie)
}
