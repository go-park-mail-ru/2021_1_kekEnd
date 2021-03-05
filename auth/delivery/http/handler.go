package http

import (
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2021_1_kekEnd/auth"
	"net/http"
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
