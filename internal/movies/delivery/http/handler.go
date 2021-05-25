package http

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/logger"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/movies"
	_const "github.com/go-park-mail-ru/2021_1_kekEnd/pkg/const"
)

type Handler struct {
	useCase movies.UseCase
	Log     *logger.Logger
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
		Log:     Log,
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
	auth, ok := ctx.Get(_const.AuthStatusKey)
	authBool := auth.(bool)
	username := ""
	if ok && authBool {
		user, ok := ctx.Get(_const.UserKey)
		if ok {
			userModel := user.(models.User)
			username = userModel.Username
		}
	}

	movie, err := h.useCase.GetMovie(ctx.Param("id"), username)
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
	} else if category == "genre" {
		h.GetMoviesByGenres(ctx)
	}
}

func (h *Handler) GetBestMovies(ctx *gin.Context) {
	auth, ok := ctx.Get(_const.AuthStatusKey)
	authBool := auth.(bool)
	username := ""
	if ok && authBool {
		user, ok := ctx.Get(_const.UserKey)
		if ok {
			userModel := user.(models.User)
			username = userModel.Username
		}
	}

	page, err := strconv.Atoi(ctx.DefaultQuery("page", _const.PageDefault))
	if err != nil || page < 1 {
		h.Log.LogWarning(ctx, "movie", "GetBestMovies", err.Error())
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	pagesNumber, bestMovies, err := h.useCase.GetBestMovies(page, username)

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

func (h *Handler) GetGenres(ctx *gin.Context) {
	genres, err := h.useCase.GetAllGenres()
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	ctx.JSON(http.StatusOK, genres)
}

func (h *Handler) GetMoviesByGenres(ctx *gin.Context) {
	auth, ok := ctx.Get(_const.AuthStatusKey)
	authBool := auth.(bool)
	username := ""
	if ok && authBool {
		user, ok := ctx.Get(_const.UserKey)
		if ok {
			userModel := user.(models.User)
			username = userModel.Username
		}
	}

	genresQuery := ctx.Query("filter")
	if genresQuery == "" {
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	genres := strings.Split(genresQuery, " ")

	page, err := strconv.Atoi(ctx.DefaultQuery("page", _const.PageDefault))
	if err != nil || page < 1 {
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	pagesNumber, moviesList, err := h.useCase.GetMoviesByGenres(genres, page, username)

	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	moviesResponse := moviesPageResponse{
		CurrentPage: page,
		PagesNumber: pagesNumber,
		MaxItems:    _const.MoviesPageSize,
		Movies:      moviesList,
	}

	ctx.JSON(http.StatusOK, moviesResponse)
}

func (h *Handler) MarkWatched(ctx *gin.Context) {
	user, ok := ctx.Get(_const.UserKey)
	if !ok {
		err := fmt.Errorf("%s", "Failed to retrieve user from context")
		h.Log.LogWarning(ctx, "movies", "MarkWatched", err.Error())
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	userModel, ok := user.(models.User)
	if !ok {
		err := fmt.Errorf("%s", "Failed to cast user to model")
		h.Log.LogError(ctx, "movies", "MarkWatched", err)
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	id := ctx.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		h.Log.LogError(ctx, "movies", "MarkWatched", err)
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}
	err = h.useCase.MarkWatched(userModel, idInt)
	if err != nil {
		h.Log.LogError(ctx, "movies", "MarkWatched", err)
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	ctx.Status(http.StatusOK)
}

func (h *Handler) MarkUnwatched(ctx *gin.Context) {
	user, ok := ctx.Get(_const.UserKey)
	if !ok {
		err := fmt.Errorf("%s", "Failed to retrieve user from context")
		h.Log.LogWarning(ctx, "movies", "MarkUnwatched", err.Error())
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	userModel, ok := user.(models.User)
	if !ok {
		err := fmt.Errorf("%s", "Failed to cast user to model")
		h.Log.LogError(ctx, "movies", "MarkUnwatched", err)
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	id := ctx.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		h.Log.LogError(ctx, "movies", "MarkUnwatched", err)
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}
	err = h.useCase.MarkUnwatched(userModel, idInt)
	if err != nil {
		h.Log.LogError(ctx, "movies", "MarkUnwatched", err)
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	ctx.Status(http.StatusOK)
}

func (h *Handler) GetSimilar(ctx *gin.Context) {
	similarMovies, err := h.useCase.GetSimilar(ctx.Param("id"))
	if err != nil {
		h.Log.LogError(ctx, "movies", "GetSimilar", err)
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	ctx.JSON(http.StatusOK, similarMovies)
}
