package http

import (
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/sessions"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/users"
	"net/http"
	"time"
)

type Handler struct {
	useCase  users.UseCase
	sessions sessions.Delivery
}

func NewHandler(useCase users.UseCase, sessions sessions.Delivery) *Handler {
	return &Handler{
		useCase:  useCase,
		sessions: sessions,
	}
}

type signupData struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) CreateUser(ctx *gin.Context) {
	signupData := new(signupData)

	err := ctx.BindJSON(signupData)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
	}

	user := &models.User{
		Username:      signupData.Username,
		Email:         signupData.Email,
		Password:      signupData.Password,
		MoviesWatched: 0,
		ReviewsNumber: 0,
	}

	err = h.useCase.CreateUser(user)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
	}

	//refactor it later
	expires := 240 * time.Hour
	userSessionID, err := h.sessions.Create(ctx, signupData.Username, expires)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
	}

	ctx.SetCookie(
		"session_id",
		userSessionID,
		int(expires),
		"/",
		ctx.Request.Host, //???
		true,
		true,
	)

	ctx.Status(http.StatusCreated) // 201
}

type loginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *Handler) Login(ctx *gin.Context) {
	loginData := new(loginData)

	err := ctx.BindJSON(loginData)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
	}

	loginStatus := h.useCase.Login(loginData.Username, loginData.Password)
	if !loginStatus {
		ctx.AbortWithStatus(http.StatusUnauthorized) // 401
	}

	//refactor it later
	expires := 240 * time.Hour
	userSessionID, err := h.sessions.Create(ctx, loginData.Username, expires)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
	}

	ctx.SetCookie(
		"session_id",
		userSessionID,
		int(expires),
		"/",
		ctx.Request.Host, //???
		false,
		true,
	)

	ctx.Status(http.StatusOK) // 200
}

func (h *Handler) GetUser(ctx *gin.Context) {
	user, err := h.useCase.GetUser(ctx.Param("username"))
	if err != nil {
		ctx.AbortWithStatus(http.StatusNotFound) // 404
	}

	ctx.JSON(http.StatusOK, user)
}

func (h *Handler) UpdateUser(ctx *gin.Context) {
	user := new(models.User)
	err := ctx.BindJSON(user)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
	}

	err = h.useCase.UpdateUser(ctx.Param("username"), user)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
	}

	ctx.JSON(http.StatusOK, user)
}
