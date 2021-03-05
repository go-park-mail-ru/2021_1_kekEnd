package http

import (
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2021_1_kekEnd/auth"
	"net/http"
	"strconv"
)

type Handler struct {
	useCase auth.UseCase
}

func NewHandler(useCase auth.UseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

type signupData struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) SignUp(ctx *gin.Context) {
	signupData := new(signupData)

	if err := ctx.BindJSON(signupData); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	if err := h.useCase.SignUp(ctx.Request.Context(), signupData.Username, signupData.Email, signupData.Password); err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
	}

	ctx.Status(http.StatusOK) // 200
}

type loginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *Handler) Login(ctx *gin.Context) {
	loginData := new(loginData)

	if err := ctx.BindJSON(loginData); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	loginStatus := h.useCase.Login(ctx.Request.Context(), loginData.Username, loginData.Password)
	if !loginStatus {
		ctx.AbortWithStatus(http.StatusUnauthorized) // 401
		return
	}

	ctx.Status(http.StatusOK) // 200
}

func (h *Handler) GetUser(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user, err := h.useCase.GetUser(ctx.Request.Context(), id)
	if err != nil {
		ctx.AbortWithStatus(http.StatusNotFound) // 404
		return
	}

	ctx.JSON(http.StatusOK, user)
}
