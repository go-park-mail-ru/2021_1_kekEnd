package http

import (
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/movies"
	_const "github.com/go-park-mail-ru/2021_1_kekEnd/pkg/const"
	"net/http"
	"strconv"
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
		return
	}

	err = h.useCase.CreateMovie(movieData)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	ctx.Status(http.StatusCreated)
}

func (h *Handler) GetMovie(ctx *gin.Context) {
	movie, err := h.useCase.GetMovie(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatus(http.StatusNotFound) // 404
		return
	}

	ctx.JSON(http.StatusOK, movie)
}

func (h *Handler) GetBestMovies(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.DefaultQuery("page", _const.PageDefault))
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	pagesNumber, bestMovies := h.useCase.GetBestMovies(page)

	ctx.JSON(http.StatusOK, gin.H{
		"current_page": page,
		"pages_number": pagesNumber,
		"best_movies":  bestMovies,
	})
}
