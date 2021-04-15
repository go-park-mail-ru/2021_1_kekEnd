package http

import (
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/logger"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/movies"
	_const "github.com/go-park-mail-ru/2021_1_kekEnd/pkg/const"
	"net/http"
	"strconv"
)

type Handler struct {
	useCase movies.UseCase
	Log *logger.Logger
}

type moviesPageResponse struct {
	CurrentPage int             `json:"current_page"`
	PagesNumber int             `json:"pages_number"`
	MaxItems    int             `json:"max_items"`
	Movies      []*models.Movie `json:"movies"`
}

func NewHandler(useCase movies.UseCase, Log *logger.Logger) *Handler {
	return &Handler{
		useCase: useCase,
		Log: Log,
	}
}

func (h *Handler) CreateMovie(ctx *gin.Context) {
	movieData := new(models.Movie)
	err := ctx.BindJSON(movieData)
	if err != nil {
		h.Log.LogWarning(ctx, "movie", "CreateMovie", err.Error())
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	err = h.useCase.CreateMovie(movieData)
	if err != nil {
		h.Log.LogError(ctx, "movie", "CreateMovie", err)
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	ctx.Status(http.StatusCreated)
}

func (h *Handler) GetMovie(ctx *gin.Context) {
	movie, err := h.useCase.GetMovie(ctx.Param("id"))
	if err != nil {
		h.Log.LogWarning(ctx, "movie", "GetMovie", err.Error())
		ctx.AbortWithStatus(http.StatusNotFound) // 404
		return
	}

	ctx.JSON(http.StatusOK, movie)
}

func (h *Handler) GetMovies(ctx *gin.Context) {
	category := ctx.Query("category")
	if category == "best" {
		h.GetBestMovies(ctx)
	}
}

func (h *Handler) GetBestMovies(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.DefaultQuery("page", _const.PageDefault))
	if err != nil || page < 1 {
		h.Log.LogWarning(ctx, "movie", "GetBestMovies", err.Error())
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	pagesNumber, bestMovies, err := h.useCase.GetBestMovies(page)

	if err != nil {
		h.Log.LogError(ctx, "movie", "GetBestMovies", err)
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	moviesResponse := moviesPageResponse{
		CurrentPage: page,
		PagesNumber: pagesNumber,
		MaxItems:    _const.MoviesPageSize,
		Movies:      bestMovies,
	}

	ctx.JSON(http.StatusOK, moviesResponse)
}
