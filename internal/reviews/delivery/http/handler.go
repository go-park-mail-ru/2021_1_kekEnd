package http

import (
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/reviews"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/users"
	_const "github.com/go-park-mail-ru/2021_1_kekEnd/pkg/const"
	"net/http"
	"strconv"
)

type Handler struct {
	reviewsUC reviews.UseCase
	usersUC   users.UseCase
}

type ReviewsResponse struct {
	CurrentPage int              `json:"current_page"`
	PagesNumber int              `json:"pages_number"`
	Movies      []*models.Review `json:"reviews"`
}

func NewHandler(useCase reviews.UseCase, usersUC users.UseCase) *Handler {
	return &Handler{
		reviewsUC: useCase,
		usersUC:   usersUC,
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

	err = h.reviewsUC.CreateReview(&userModel, review)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
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

	userReviews, err := h.reviewsUC.GetReviewsByUser(userModel.Username)

	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	ctx.JSON(http.StatusOK, userReviews)
}

func (h *Handler) GetMovieReviews(ctx *gin.Context) {
	movieID := ctx.Param("id")
	page, err := strconv.Atoi(ctx.DefaultQuery("page", _const.PageDefault))
	if err != nil || page < 1 {
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	pagesNumber, movieReviews, err := h.reviewsUC.GetReviewsByMovie(movieID, page)

	if err != nil {
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

func (h *Handler) EditUserReviewForMovie(ctx *gin.Context) {
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

	err = h.reviewsUC.EditUserReviewForMovie(&userModel, review)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	ctx.Status(http.StatusOK) // 200
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

	err := h.reviewsUC.DeleteUserReviewForMovie(&userModel, movieID)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	ctx.Status(http.StatusOK) // 200
}
