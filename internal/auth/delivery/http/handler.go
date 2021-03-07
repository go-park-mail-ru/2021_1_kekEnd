package http

import (
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/auth"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
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
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
}

func (h *Handler) SignUp(ctx *gin.Context) {
	signupData := new(signupData)

	if err := ctx.BindJSON(signupData); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
	}

	if err := h.useCase.SignUp(signupData.Username, signupData.Email, signupData.FirstName,
		signupData.LastName, signupData.Password); err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
	}

	ctx.Status(http.StatusCreated) // 201
}

type loginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *Handler) Login(ctx *gin.Context) {
	loginData := new(loginData)

	if err := ctx.BindJSON(loginData); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
	}

	loginStatus := h.useCase.Login(loginData.Username, loginData.Password)
	if !loginStatus {
		ctx.AbortWithStatus(http.StatusUnauthorized) // 401
	}

	ctx.Status(http.StatusOK) // 200
}

func (h *Handler) GetUser(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
	}

	user, err := h.useCase.GetUser(id)
	if err != nil {
		ctx.AbortWithStatus(http.StatusNotFound) // 404
	}

	ctx.JSON(http.StatusOK, user)
}

func (h *Handler) UpdateUser(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
	}

	user := new(models.User)
	if err := ctx.BindJSON(user); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
	}

	if err := h.useCase.UpdateUser(id, user); err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
	}

	// TODO: отправлять либо 200, либо 201
	ctx.JSON(http.StatusOK, user)
}
