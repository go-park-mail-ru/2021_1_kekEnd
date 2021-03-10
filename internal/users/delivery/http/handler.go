package http

import (
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/sessions"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/users"
	"github.com/google/uuid"
	"net/http"
	"path/filepath"
	"time"
)

const userKey = "user"
const host = "localhost"

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
		return
	}

	if signupData.Username == "" || signupData.Email == "" || signupData.Password == ""{
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	//path, _ := filepath.Abs("tmp/avatars/default.jpeg")
	defaultAvatar := "tmp/avatars/default.jpeg"
	user := &models.User{
		Username:      signupData.Username,
		Email:         signupData.Email,
		Password:      signupData.Password,
		Avatar:		   defaultAvatar,
		MoviesWatched: 0,
		ReviewsNumber: 0,
	}

	err = h.useCase.CreateUser(user)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	//refactor it later
	expires := 240 * time.Hour
	userSessionID, err := h.sessions.Create(ctx, signupData.Username, expires)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	ctx.SetCookie(
		"session_id",
		userSessionID,
		int(expires),
		"/",
		host,
		false,
		true,
	)

	ctx.Status(http.StatusCreated) // 201
}

type loginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *Handler) Logout(ctx *gin.Context) {
	cookie, err := ctx.Cookie("session_id")
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized) // 401
		return
	}

	err = h.sessions.Delete(ctx, cookie)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	ctx.SetCookie("session_id", "Delete cookie", -1, "/", host, false, true)

	ctx.Status(http.StatusOK) // 200
}

func (h *Handler) Login(ctx *gin.Context) {
	loginData := new(loginData)

	err := ctx.BindJSON(loginData)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	loginStatus := h.useCase.Login(loginData.Username, loginData.Password)
	if !loginStatus {
		ctx.AbortWithStatus(http.StatusUnauthorized) // 401
		return
	}

	//refactor it later
	expires := 240 * time.Hour
	userSessionID, err := h.sessions.Create(ctx, loginData.Username, expires)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	ctx.SetCookie(
		"session_id",
		userSessionID,
		int(expires),
		"/",
		host,
		false,
		true,
	)

	ctx.Status(http.StatusOK) // 200
}

func (h *Handler) GetUser(ctx *gin.Context) {
	user, ok := ctx.Get(userKey)
	if !ok {
		ctx.AbortWithStatus(http.StatusNotFound) // 404
		return
	}

	userModel, ok := user.(models.User)
	if !ok {
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	userNoPassword := models.FromUser(userModel)
	ctx.JSON(http.StatusOK, userNoPassword)
}

func (h *Handler) UpdateUser(ctx *gin.Context) {
	changed := new(models.User)
	err := ctx.BindJSON(changed)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	user, ok := ctx.Get(userKey)
	if !ok {
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	userModel, ok := user.(models.User)
	if !ok {
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	newUser, err := h.useCase.UpdateUser(&userModel, *changed)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	userNoPassword := models.FromUser(*newUser)
	ctx.JSON(http.StatusOK, userNoPassword)
}

func (h *Handler) UploadAvatar(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	extension := filepath.Ext(file.Filename)
	// generate random file name for the new uploaded file so it doesn't override the old file with same name
	newFileName := uuid.New().String() + extension
	err = ctx.SaveUploadedFile(file, "tmp/avatars/" + newFileName)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	// TODO: add avatar reference to user model. Set it here
	user, ok := ctx.Get(userKey)
	if !ok {
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	userModel := user.(models.User)
	change := userModel

	//avatar := make([]byte, file.Size)
	//fileReader, _ := file.Open()
	//_, err = fileReader.Read(avatar)
	//if err != nil {
	//	ctx.AbortWithStatus(http.StatusInternalServerError) // 500
	//	return
	//}

	change.Avatar = "tmp/avatars/" + newFileName
	_, err = h.useCase.UpdateUser(&userModel, change)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	ctx.Status(http.StatusOK)
}
