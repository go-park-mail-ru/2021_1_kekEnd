package http

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/logger"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/playlists"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/users"
	constants "github.com/go-park-mail-ru/2021_1_kekEnd/pkg/const"
)

// Handler структура хендлера
type Handler struct {
	useCase      playlists.UseCase
	usersUseCase users.UseCase
	Log          *logger.Logger
}

// PlaylistMovie структура фильма для плейлиста
type PlaylistMovie struct {
	MovieID string `json:"movie_id"`
}

// PlaylistUser структура фильма для юзера
type PlaylistUser struct {
	Username string `json:"username"`
}

// NewHandler инициализация хендлера
func NewHandler(useCase playlists.UseCase, usersUseCase users.UseCase, Log *logger.Logger) *Handler {
	return &Handler{
		useCase:      useCase,
		usersUseCase: usersUseCase,
		Log:          Log,
	}
}

// CreatePlaylist создание плейлиста
func (h *Handler) CreatePlaylist(ctx *gin.Context) {
	playlistData := new(models.Playlist)
	err := ctx.BindJSON(playlistData)
	if err != nil {
		msg := "Failed to bind request data " + err.Error()
		h.Log.LogWarning(ctx, "playlists", "CreatePlaylist", msg)
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	user, ok := ctx.Get(constants.UserKey)
	if !ok {
		err := fmt.Errorf("%s", "Failed to retrieve user from context")
		h.Log.LogWarning(ctx, "playlists", "CreatePlaylist", err.Error())
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	userModel, ok := user.(models.User)
	if !ok {
		err := fmt.Errorf("%s", "Failed to cast user to model")
		h.Log.LogError(ctx, "playlists", "CreatePlaylist", err)
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
	}

	err = h.useCase.CreatePlaylist(userModel.Username, playlistData.Name, playlistData.IsShared)
	if err != nil {
		h.Log.LogError(ctx, "playlists", "CreatePlaylist", err)
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	ctx.Status(http.StatusCreated)
}

// GetPlaylist получение плейлиста
func (h *Handler) GetPlaylist(ctx *gin.Context) {
	playlistIDStr := ctx.Param("playlist_id")

	playlistID, err := strconv.Atoi(playlistIDStr)
	if err != nil {
		err := fmt.Errorf("%s", "Failed to cast playlistID value to number")
		h.Log.LogWarning(ctx, "playlists", "GetPlaylist", err.Error())
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	playlist, err := h.useCase.GetPlaylist(playlistID)
	if err != nil {
		h.Log.LogWarning(ctx, "playlists", "GetPlaylist", err.Error())
		ctx.AbortWithStatus(http.StatusNotFound) // 404
		return
	}

	ctx.JSON(http.StatusOK, playlist)
}

// GetPlaylistsInfo получение информации о плейлисте
func (h *Handler) GetPlaylistsInfo(ctx *gin.Context) {
	movieIDStr := ctx.Param("movie_id")
	user, ok := ctx.Get(constants.UserKey)
	if !ok {
		err := fmt.Errorf("%s", "Failed to retrieve user from context")
		h.Log.LogWarning(ctx, "playlists", "GetPlaylistsInfo", err.Error())
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	userModel, ok := user.(models.User)
	if !ok {
		err := fmt.Errorf("%s", "Failed to cast user to model")
		h.Log.LogError(ctx, "playlists", "GetPlaylistsInfo", err)
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	movieID, err := strconv.Atoi(movieIDStr)
	if err != nil {
		err := fmt.Errorf("%s", "Failed to cast movieID value to number")
		h.Log.LogWarning(ctx, "playlists", "GetPlaylistsInfo", err.Error())
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	playlistInfo, err := h.useCase.GetPlaylistsInfo(userModel.Username, movieID)
	if err != nil {
		h.Log.LogWarning(ctx, "playlists", "GetPlaylistsInfo", err.Error())
		ctx.AbortWithStatus(http.StatusNotFound) // 404
		return
	}

	ctx.JSON(http.StatusOK, playlistInfo)
}

// GetPlaylists получить все плейлисты
func (h *Handler) GetPlaylists(ctx *gin.Context) {
	userModel, err := h.usersUseCase.GetUser(ctx.Param("username"))
	if err != nil {
		err := fmt.Errorf("%s", "Failed to get user")
		h.Log.LogError(ctx, "playlists", "GetPlaylists", err)
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	userPlaylists, err := h.useCase.GetPlaylists(userModel.Username)
	if err != nil {
		h.Log.LogWarning(ctx, "playlists", "GetPlaylists", err.Error())
		ctx.AbortWithStatus(http.StatusNotFound) // 404
		return
	}

	ctx.JSON(http.StatusOK, userPlaylists)
}

// EditPlaylist изменить плейлист
func (h *Handler) EditPlaylist(ctx *gin.Context) {
	playlistData := new(models.Playlist)
	err := ctx.BindJSON(playlistData)
	if err != nil {
		msg := "Failed to bind request data " + err.Error()
		h.Log.LogWarning(ctx, "playlists", "EditPlaylist", msg)
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	user, ok := ctx.Get(constants.UserKey)
	if !ok {
		err := fmt.Errorf("%s", "Failed to retrieve user from context")
		h.Log.LogWarning(ctx, "playlists", "EditPlaylist", err.Error())
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	userModel, ok := user.(models.User)
	if !ok {
		err := fmt.Errorf("%s", "Failed to cast user to model")
		h.Log.LogError(ctx, "playlists", "EditPlaylist", err)
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	playlistID, err := strconv.Atoi(playlistData.ID)
	if err != nil {
		err := fmt.Errorf("%s", "Failed to cast playlistID value to number")
		h.Log.LogWarning(ctx, "playlists", "EditPlaylist", err.Error())
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	err = h.useCase.UpdatePlaylist(userModel.Username, playlistID, playlistData.Name, playlistData.IsShared)
	if err != nil {
		h.Log.LogError(ctx, "playlists", "EditPlaylist", err)
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	ctx.Status(http.StatusOK)
}

// DeletePlaylist удалить плейлист
func (h *Handler) DeletePlaylist(ctx *gin.Context) {
	playlistIDStr := ctx.Param("playlist_id")
	user, ok := ctx.Get(constants.UserKey)
	if !ok {
		err := fmt.Errorf("%s", "Failed to retrieve user from context")
		h.Log.LogWarning(ctx, "playlists", "DeletePlaylist", err.Error())
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	userModel, ok := user.(models.User)
	if !ok {
		err := fmt.Errorf("%s", "Failed to cast user to model")
		h.Log.LogError(ctx, "playlists", "DeletePlaylist", err)
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	playlistID, err := strconv.Atoi(playlistIDStr)
	if err != nil {
		err := fmt.Errorf("%s", "Failed to cast playlistID value to number")
		h.Log.LogWarning(ctx, "playlists", "DeletePlaylist", err.Error())
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	err = h.useCase.DeletePlaylist(userModel.Username, playlistID)
	if err != nil {
		h.Log.LogWarning(ctx, "playlists", "DeletePlaylist", err.Error())
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	ctx.Status(http.StatusOK)
}

// AddMovieToPlaylist добавить фильм в плейлист
func (h *Handler) AddMovieToPlaylist(ctx *gin.Context) {
	playlistIDStr := ctx.Param("playlist_id")
	playlistMovieData := new(PlaylistMovie)
	err := ctx.BindJSON(playlistMovieData)
	if err != nil {
		msg := "Failed to bind request data " + err.Error()
		h.Log.LogWarning(ctx, "playlists", "AddMovieToPlaylist", msg)
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	user, ok := ctx.Get(constants.UserKey)
	if !ok {
		err := fmt.Errorf("%s", "Failed to retrieve user from context")
		h.Log.LogWarning(ctx, "playlists", "AddMovieToPlaylist", err.Error())
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	userModel, ok := user.(models.User)
	if !ok {
		err := fmt.Errorf("%s", "Failed to cast user to model")
		h.Log.LogError(ctx, "playlists", "AddMovieToPlaylist", err)
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
	}

	playlistID, err := strconv.Atoi(playlistIDStr)
	if err != nil {
		err := fmt.Errorf("%s", "Failed to cast playlistID value to number")
		h.Log.LogWarning(ctx, "playlists", "DeletePlaylist", err.Error())
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	movieID, err := strconv.Atoi(playlistMovieData.MovieID)
	if err != nil {
		err := fmt.Errorf("%s", "Failed to cast movieID value to number")
		h.Log.LogWarning(ctx, "playlists", "AddMovieToPlaylist", err.Error())
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	err = h.useCase.AddMovieToPlaylist(userModel.Username, playlistID, movieID)
	if err != nil {
		h.Log.LogError(ctx, "playlists", "AddMovieToPlaylist", err)
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	ctx.Status(http.StatusCreated)
}

// DeleteMovieFromPlaylist удалить фильм из плейлиста
func (h *Handler) DeleteMovieFromPlaylist(ctx *gin.Context) {
	playlistIDStr := ctx.Param("playlist_id")
	playlistMovieData := new(PlaylistMovie)
	err := ctx.BindJSON(playlistMovieData)
	if err != nil {
		msg := "Failed to bind request data " + err.Error()
		h.Log.LogWarning(ctx, "playlists", "DeleteMovieFromPlaylist", msg)
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	user, ok := ctx.Get(constants.UserKey)
	if !ok {
		err := fmt.Errorf("%s", "Failed to retrieve user from context")
		h.Log.LogWarning(ctx, "playlists", "DeleteMovieFromPlaylist", err.Error())
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	userModel, ok := user.(models.User)
	if !ok {
		err := fmt.Errorf("%s", "Failed to cast user to model")
		h.Log.LogError(ctx, "playlists", "DeleteMovieFromPlaylist", err)
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
	}

	playlistID, err := strconv.Atoi(playlistIDStr)
	if err != nil {
		err := fmt.Errorf("%s", "Failed to cast playlistID value to number")
		h.Log.LogWarning(ctx, "playlists", "DeletePlaylist", err.Error())
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	movieID, err := strconv.Atoi(playlistMovieData.MovieID)
	if err != nil {
		err := fmt.Errorf("%s", "Failed to cast movieID value to number")
		h.Log.LogWarning(ctx, "playlists", "DeleteMovieFromPlaylist", err.Error())
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	err = h.useCase.DeleteMovieFromPlaylist(userModel.Username, playlistID, movieID)
	if err != nil {
		h.Log.LogError(ctx, "playlists", "DeleteMovieFromPlaylist", err)
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	ctx.Status(http.StatusOK)
}

// AddUserToPlaylist добавить юзера в плейлист
func (h *Handler) AddUserToPlaylist(ctx *gin.Context) {
	playlistIDStr := ctx.Param("playlist_id")
	playlistUserData := new(PlaylistUser)
	err := ctx.BindJSON(playlistUserData)
	if err != nil {
		msg := "Failed to bind request data " + err.Error()
		h.Log.LogWarning(ctx, "playlists", "AddUserToPlaylist", msg)
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	user, ok := ctx.Get(constants.UserKey)
	if !ok {
		err := fmt.Errorf("%s", "Failed to retrieve user from context")
		h.Log.LogWarning(ctx, "playlists", "AddUserToPlaylist", err.Error())
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	userModel, ok := user.(models.User)
	if !ok {
		err := fmt.Errorf("%s", "Failed to cast user to model")
		h.Log.LogError(ctx, "playlists", "AddUserToPlaylist", err)
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
	}

	playlistID, err := strconv.Atoi(playlistIDStr)
	if err != nil {
		err := fmt.Errorf("%s", "Failed to cast playlistID value to number")
		h.Log.LogWarning(ctx, "playlists", "DeletePlaylist", err.Error())
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	err = h.useCase.AddUserToPlaylist(userModel.Username, playlistID, playlistUserData.Username)
	if err != nil {
		h.Log.LogError(ctx, "playlists", "AddUserToPlaylist", err)
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	ctx.Status(http.StatusCreated)
}

// DeleteUserFromPlaylist удалить юзера из плейлиста
func (h *Handler) DeleteUserFromPlaylist(ctx *gin.Context) {
	playlistIDStr := ctx.Param("playlist_id")
	playlistUserData := new(PlaylistUser)
	err := ctx.BindJSON(playlistUserData)
	if err != nil {
		msg := "Failed to bind request data " + err.Error()
		h.Log.LogWarning(ctx, "playlists", "DeleteUserFromPlaylist", msg)
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	user, ok := ctx.Get(constants.UserKey)
	if !ok {
		err := fmt.Errorf("%s", "Failed to retrieve user from context")
		h.Log.LogWarning(ctx, "playlists", "DeleteUserFromPlaylist", err.Error())
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	userModel, ok := user.(models.User)
	if !ok {
		err := fmt.Errorf("%s", "Failed to cast user to model")
		h.Log.LogError(ctx, "playlists", "DeleteUserFromPlaylist", err)
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
	}

	playlistID, err := strconv.Atoi(playlistIDStr)
	if err != nil {
		err := fmt.Errorf("%s", "Failed to cast playlistID value to number")
		h.Log.LogWarning(ctx, "playlists", "DeletePlaylist", err.Error())
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	err = h.useCase.DeleteUserFromPlaylist(userModel.Username, playlistID, playlistUserData.Username)
	if err != nil {
		h.Log.LogError(ctx, "playlists", "DeleteUserFromPlaylist", err)
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	ctx.Status(http.StatusOK)
}
