package http

import (
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/reviews"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/users"
	_const "github.com/go-park-mail-ru/2021_1_kekEnd/pkg/const"
	"net/http"
)

type Handler struct {
	reviewsUC reviews.UseCase
	usersUC   users.UseCase
}

func NewHandler(useCase reviews.UseCase, usersUC users.UseCase) *Handler {
	return &Handler{
		reviewsUC: useCase,
		usersUC: usersUC,
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

	err = h.reviewsUC.CreateReview(review)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	gotUser, err := h.usersUC.GetUser(userModel.Username)
	if err != nil {
		ctx.Status(http.StatusNotFound) // 404
		return
	}
	_, err = h.usersUC.UpdateUser(gotUser, models.User{
		Username: gotUser.Username,
		ReviewsNumber: gotUser.ReviewsNumber + 1,
	})
	if err != nil {
		ctx.Status(http.StatusInternalServerError) // 500
		return
	}

	ctx.Status(http.StatusCreated) // 201
}

func (h *Handler) GetUserReviews(ctx *gin.Context) {
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

	userReviews := h.reviewsUC.GetReviewsByUser(userModel.Username)

	ctx.JSON(http.StatusOK, userReviews)
}

func (h *Handler) GetMovieReviews(ctx *gin.Context) {
	movieID := ctx.Param("id")
	movieReviews := h.reviewsUC.GetReviewsByMovie(movieID)
	ctx.JSON(http.StatusOK, movieReviews)
}

func (h *Handler) GetUserReviewForMovie(ctx *gin.Context) {
	movieID := ctx.Param("id")

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

	review, err := h.reviewsUC.GetUserReviewForMovie(userModel.Username, movieID)
	if err != nil {
		ctx.AbortWithStatus(http.StatusNotFound) // 404
		return
	}

	ctx.JSON(http.StatusOK, review)
}

func (h *Handler) DeleteUserReviewForMovie(ctx *gin.Context) {
	movieID := ctx.Param("id")

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

	err := h.reviewsUC.DeleteUserReviewForMovie(userModel.Username, movieID)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	ctx.Status(http.StatusOK) // 200
}
