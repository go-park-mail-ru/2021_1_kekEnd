package http

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/logger"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/reviews"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/users"
	constants "github.com/go-park-mail-ru/2021_1_kekEnd/pkg/const"
)

// Handler структура хендлера
type Handler struct {
	reviewsUC reviews.UseCase
	usersUC   users.UseCase
	Log       *logger.Logger
}

// ReviewsResponse структура рецензии
type ReviewsResponse struct {
	CurrentPage int              `json:"current_page"`
	PagesNumber int              `json:"pages_number"`
	Movies      []*models.Review `json:"reviews"`
}

// NewHandler инициализация хендлера рецензий
func NewHandler(useCase reviews.UseCase, usersUC users.UseCase, Log *logger.Logger) *Handler {
	return &Handler{
		reviewsUC: useCase,
		usersUC:   usersUC,
		Log:       Log,
	}
}

// CreateReview создание рецензии
func (h *Handler) CreateReview(ctx *gin.Context) {
	review := new(models.Review)
	err := ctx.BindJSON(review)
	if err != nil {
		msg := "Failed to bind request data" + err.Error()
		h.Log.LogWarning(ctx, "reviews", "CreateReview", msg)
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	user, ok := ctx.Get(constants.UserKey)
	if !ok {
		err := fmt.Errorf("%s", "Failed to retrieve user from context")
		h.Log.LogError(ctx, "reviews", "CreateReview", err)
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	userModel, ok := user.(models.User)
	if !ok {
		err := fmt.Errorf("%s", "Failed to cast user to model")
		h.Log.LogError(ctx, "reviews", "CreateReview", err)
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	err = h.reviewsUC.CreateReview(&userModel, review)
	if err != nil {
		h.Log.LogWarning(ctx, "reviews", "CreateReview", err.Error())
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	ctx.Status(http.StatusCreated) // 201
}

// GetUserReviews полученить рецензии пользователя
func (h *Handler) GetUserReviews(ctx *gin.Context) {
	userModel, err := h.usersUC.GetUser(ctx.Param("username"))
	if err != nil {
		err := fmt.Errorf("%s", "Failed to get user")
		h.Log.LogError(ctx, "reviews", "GetUserReviews", err)
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	userReviews, err := h.reviewsUC.GetReviewsByUser(userModel.Username)

	if err != nil {
		h.Log.LogError(ctx, "reviews", "GetUserReviews", err)
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	ctx.JSON(http.StatusOK, userReviews)
}

// GetMovieReviews получить рецензии фильма
func (h *Handler) GetMovieReviews(ctx *gin.Context) {
	movieID := ctx.Param("id")
	page, err := strconv.Atoi(ctx.DefaultQuery("page", constants.PageDefault))
	if err != nil || page < 1 {
		err := fmt.Errorf("%s", "Failed to cast 'page' to number or invalid page")
		h.Log.LogWarning(ctx, "reviews", "GetMovieReviews", err.Error())
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	pagesNumber, movieReviews, err := h.reviewsUC.GetReviewsByMovie(movieID, page)

	if err != nil {
		h.Log.LogError(ctx, "reviews", "GetMovieReviews", err)
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	reviewsResponse := ReviewsResponse{
		CurrentPage: page,
		PagesNumber: pagesNumber,
		Movies:      movieReviews,
	}

	ctx.JSON(http.StatusOK, reviewsResponse)
}

// GetUserReviewForMovie получить рецензию пользователя к фильму
func (h *Handler) GetUserReviewForMovie(ctx *gin.Context) {
	movieID := ctx.Param("id")

	user, ok := ctx.Get(constants.UserKey)
	if !ok {
		err := fmt.Errorf("%s", "Failed to retrieve user from context")
		h.Log.Error(ctx, "reviews", "GetUserReviewForMovie", err.Error)
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	userModel, ok := user.(models.User)
	if !ok {
		err := fmt.Errorf("%s", "Failed to cast user to model")
		h.Log.Error(ctx, "reviews", "GetUserReviewForMovie", err.Error)
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	review, err := h.reviewsUC.GetUserReviewForMovie(userModel.Username, movieID)
	if err != nil {
		h.Log.LogWarning(ctx, "reviews", "GetUserReviewForMovie", err.Error())
		ctx.AbortWithStatus(http.StatusNotFound) // 404
		return
	}

	ctx.JSON(http.StatusOK, review)
}

// EditUserReviewForMovie изменить рецензию пользователя
func (h *Handler) EditUserReviewForMovie(ctx *gin.Context) {
	review := new(models.Review)
	err := ctx.BindJSON(review)
	if err != nil {
		msg := "Failed to bind request data" + err.Error()
		h.Log.LogWarning(ctx, "reviews", "EditUserReviewForMovie", msg)
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	user, ok := ctx.Get(constants.UserKey)
	if !ok {
		err := fmt.Errorf("%s", "Failed to retrieve user from context")
		h.Log.LogError(ctx, "reviews", "EditUserReviewForMovie", err)
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	userModel, ok := user.(models.User)
	if !ok {
		err := fmt.Errorf("%s", "Failed to cast user to model")
		h.Log.LogError(ctx, "reviews", "EditUserReviewForMovie", err)
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	err = h.reviewsUC.EditUserReviewForMovie(&userModel, review)
	if err != nil {
		h.Log.LogError(ctx, "reviews", "EditUserReviewForMovie", err)
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	ctx.Status(http.StatusOK) // 200
}

// DeleteUserReviewForMovie удалить рецензнию пользователя
func (h *Handler) DeleteUserReviewForMovie(ctx *gin.Context) {
	movieID := ctx.Param("id")

	user, ok := ctx.Get(constants.UserKey)
	if !ok {
		err := fmt.Errorf("%s", "Failed to retrieve user from context")
		h.Log.LogError(ctx, "reviews", "DeleteUserReviewForMovie", err)
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	userModel, ok := user.(models.User)
	if !ok {
		err := fmt.Errorf("%s", "Failed to cast user to model")
		h.Log.LogError(ctx, "reviews", "DeleteUserReviewForMovie", err)
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	err := h.reviewsUC.DeleteUserReviewForMovie(&userModel, movieID)
	if err != nil {
		h.Log.LogError(ctx, "reviews", "DeleteUserReviewForMovie", err)
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	ctx.Status(http.StatusOK) // 200
}

// DeleteReview удалить рецензнию пользователя
func (h *Handler) DeleteReview(ctx *gin.Context) {
	movieID := ctx.Param("id")
	movieIDInt, _ := strconv.Atoi(movieID)

	username := ctx.Param("username")
	// user, ok := ctx.Get(constants.UserKey)
	// if !ok {
	// 	err := fmt.Errorf("%s", "Failed to retrieve user from context")
	// 	h.Log.LogError(ctx, "reviews", "DeleteUserReviewForMovie", err)
	// 	ctx.AbortWithStatus(http.StatusInternalServerError) // 500
	// 	return
	// }

	// userModel, ok := user.(models.User)
	// if !ok {
	// 	err := fmt.Errorf("%s", "Failed to cast user to model")
	// 	h.Log.LogError(ctx, "reviews", "DeleteUserReviewForMovie", err)
	// 	ctx.AbortWithStatus(http.StatusInternalServerError) // 500
	// 	return
	// }

	err := h.reviewsUC.DeleteReview("admin1", username, movieIDInt)
	if err != nil {
		h.Log.LogError(ctx, "reviews", "DeleteUserReviewForMovie", err)
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	ctx.Status(http.StatusOK) // 200
}
