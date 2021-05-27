package ratings

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/logger"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/ratings"
	constants "github.com/go-park-mail-ru/2021_1_kekEnd/pkg/const"
)

// Handler структура хендлера
type Handler struct {
	useCase ratings.UseCase
	Log     *logger.Logger
}

// NewHandler инициализация нового хендлера
func NewHandler(useCase ratings.UseCase, Log *logger.Logger) *Handler {
	return &Handler{
		useCase: useCase,
		Log:     Log,
	}
}

// ratingData структура оценок
type ratingData struct {
	MovieID string `json:"movie_id"`
	Score   string `json:"score"`
}

// CreateRating создание оценки
func (h *Handler) CreateRating(ctx *gin.Context) {
	ratingData := new(ratingData)
	err := ctx.BindJSON(ratingData)
	if err != nil {
		msg := "Failed to bind request data" + err.Error()
		h.Log.LogWarning(ctx, "ratings", "CreateRating", msg)
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	user, ok := ctx.Get(constants.UserKey)
	if !ok {
		err := fmt.Errorf("%s", "Failed to retrieve user from context")
		h.Log.LogWarning(ctx, "ratings", "CreateRating", err.Error())
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	userModel, ok := user.(models.User)
	if !ok {
		err := fmt.Errorf("%s", "Failed to cast user to model")
		h.Log.LogError(ctx, "ratings", "CreateRating", err)
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
	}

	score, err := strconv.Atoi(ratingData.Score)
	if err != nil {
		err := fmt.Errorf("%s", "Failed to cast rating value to number")
		h.Log.LogWarning(ctx, "ratings", "CreateRating", err.Error())
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	err = h.useCase.CreateRating(userModel.Username, ratingData.MovieID, score)
	if err != nil {
		h.Log.LogError(ctx, "ratings", "CreateRating", err)
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	ctx.Status(http.StatusCreated)
}

// GetRating получение оценки
func (h *Handler) GetRating(ctx *gin.Context) {
	movieID := ctx.Param("movie_id")
	user, ok := ctx.Get(constants.UserKey)
	if !ok {
		err := fmt.Errorf("%s", "Failed to retrieve user from context")
		h.Log.LogWarning(ctx, "ratings", "GetRating", err.Error())
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	userModel, ok := user.(models.User)
	if !ok {
		err := fmt.Errorf("%s", "Failed to cast user to model")
		h.Log.LogError(ctx, "ratings", "GetRating", err)
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	rating, err := h.useCase.GetRating(userModel.Username, movieID)
	if err != nil {
		h.Log.LogWarning(ctx, "ratings", "GetRating", err.Error())
		ctx.AbortWithStatus(http.StatusNotFound) // 404
		return
	}

	ctx.JSON(http.StatusOK, rating)
}

// UpdateRating обновление оценки
func (h *Handler) UpdateRating(ctx *gin.Context) {
	ratingData := new(ratingData)
	err := ctx.BindJSON(ratingData)
	if err != nil {
		msg := "Failed to bind request data" + err.Error()
		h.Log.LogWarning(ctx, "ratings", "UpdateRating", msg)
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	user, ok := ctx.Get(constants.UserKey)
	if !ok {
		err := fmt.Errorf("%s", "Failed to retrieve user from context")
		h.Log.LogWarning(ctx, "ratings", "UpdateRating", err.Error())
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	userModel, ok := user.(models.User)
	if !ok {
		err := fmt.Errorf("%s", "Failed to cast user to model")
		h.Log.LogError(ctx, "ratings", "UpdateRating", err)
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	score, err := strconv.Atoi(ratingData.Score)
	if err != nil {
		msg := "Failed to cast rating value to number" + err.Error()
		h.Log.LogWarning(ctx, "ratings", "UpdateRating", msg)
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	err = h.useCase.UpdateRating(userModel.Username, ratingData.MovieID, score)
	if err != nil {
		h.Log.LogError(ctx, "ratings", "UpdateRating", err)
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	ctx.Status(http.StatusOK)
}

// DeleteRating удаление оценки
func (h *Handler) DeleteRating(ctx *gin.Context) {
	movieID := ctx.Param("movie_id")
	user, ok := ctx.Get(constants.UserKey)
	if !ok {
		err := fmt.Errorf("%s", "Failed to retrieve user from context")
		h.Log.LogWarning(ctx, "ratings", "DeleteRating", err.Error())
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	userModel, ok := user.(models.User)
	if !ok {
		err := fmt.Errorf("%s", "Failed to cast user to model")
		h.Log.LogError(ctx, "ratings", "DeleteRating", err)
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	err := h.useCase.DeleteRating(userModel.Username, movieID)
	if err != nil {
		h.Log.LogWarning(ctx, "ratings", "DeleteRating", err.Error())
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	ctx.Status(http.StatusOK)
}
