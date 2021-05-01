package playlists

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/logger"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/playlists"
	_const "github.com/go-park-mail-ru/2021_1_kekEnd/pkg/const"
)

type Handler struct {
	useCase playlists.UseCase
	Log     *logger.Logger
}

type DeleteFromPlaylist struct {
	PlaylistID string `json:"playlistID"`
	MovieID    string `json:"movieID"`
}

func NewHandler(useCase playlists.UseCase, Log *logger.Logger) *Handler {
	return &Handler{
		useCase: useCase,
		Log:     Log,
	}
}

func (h *Handler) CreatePlaylist(ctx *gin.Context) {
	playlistData := new(models.Playlist)
	err := ctx.BindJSON(playlistData)
	if err != nil {
		msg := "Failed to bind request data" + err.Error()
		h.Log.LogWarning(ctx, "playlists", "CreatePlaylist", msg)
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	user, ok := ctx.Get(_const.UserKey)
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

	playlistID, err := strconv.Atoi(playlistData.ID)
	if err != nil {
		err := fmt.Errorf("%s", "Failed to cast rating value to number")
		h.Log.LogWarning(ctx, "ratings", "CreateRating", err.Error())
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	err = h.useCase.CreatePlaylist(userModel.Username, playlistID, playlistData.IsShared)
	if err != nil {
		h.Log.LogError(ctx, "playlists", "CreateRating", err)
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	ctx.Status(http.StatusCreated)
}

func (h *Handler) GetPlaylistsInfo(ctx *gin.Context) {
	user, ok := ctx.Get(_const.UserKey)
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

	playlist, err := h.useCase.GetPlaylistsInfo(userModel.Username)
	if err != nil {
		h.Log.LogWarning(ctx, "ratings", "GetRating", err.Error())
		ctx.AbortWithStatus(http.StatusNotFound) // 404
		return
	}

	ctx.JSON(http.StatusOK, playlist)
}

func (h *Handler) EditPlaylist(ctx *gin.Context) {
	playlistData := new(models.Playlist)
	err := ctx.BindJSON(playlistData)
	if err != nil {
		msg := "Failed to bind request data" + err.Error()
		h.Log.LogWarning(ctx, "ratings", "UpdateRating", msg)
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	user, ok := ctx.Get(_const.UserKey)
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

	playlistID, err := strconv.Atoi(playlistData.ID)
	if err != nil {
		err := fmt.Errorf("%s", "Failed to cast rating value to number")
		h.Log.LogWarning(ctx, "ratings", "CreateRating", err.Error())
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	playlist, err := h.useCase.UpdatePlaylist(userModel.Username, playlistID, playlistData.IsShared)
	if err != nil {
		h.Log.LogError(ctx, "ratings", "UpdateRating", err)
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	ctx.JSON(http.StatusOK, playlist)
}

func (h *Handler) DeletePlaylist(ctx *gin.Context) {
	playlistIDStr := ctx.Param("playlist_id")
	user, ok := ctx.Get(_const.UserKey)
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

	playlistID, err := strconv.Atoi(playlistIDStr)
	if err != nil {
		err := fmt.Errorf("%s", "Failed to cast rating value to number")
		h.Log.LogWarning(ctx, "ratings", "CreateRating", err.Error())
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	err = h.useCase.DeletePlaylist(userModel.Username, playlistID)
	if err != nil {
		h.Log.LogWarning(ctx, "ratings", "DeleteRating", err.Error())
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	ctx.Status(http.StatusOK)
}

func (h *Handler) AddToPlaylist(ctx *gin.Context) {
	playlistData := new(models.Playlist)
	err := ctx.BindJSON(playlistData)
	if err != nil {
		msg := "Failed to bind request data" + err.Error()
		h.Log.LogWarning(ctx, "ratings", "CreateRating", msg)
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	user, ok := ctx.Get(_const.UserKey)
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

	err = h.useCase.AddToPlaylist(userModel.Username, playlistData)
	if err != nil {
		h.Log.LogError(ctx, "ratings", "CreateRating", err)
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	ctx.Status(http.StatusCreated)
}

func (h *Handler) DeleteFromPlaylist(ctx *gin.Context) {
	playlistData := new(DeleteFromPlaylist)
	err := ctx.BindJSON(playlistData)
	if err != nil {
		msg := "Failed to bind request data" + err.Error()
		h.Log.LogWarning(ctx, "ratings", "CreateRating", msg)
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	user, ok := ctx.Get(_const.UserKey)
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

	playlistID, err := strconv.Atoi(playlistData.PlaylistID)
	if err != nil {
		err := fmt.Errorf("%s", "Failed to cast rating value to number")
		h.Log.LogWarning(ctx, "ratings", "CreateRating", err.Error())
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	movieID, err := strconv.Atoi(playlistData.MovieID)
	if err != nil {
		err := fmt.Errorf("%s", "Failed to cast rating value to number")
		h.Log.LogWarning(ctx, "ratings", "CreateRating", err.Error())
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	err = h.useCase.DeleteFromPlaylist(userModel.Username, playlistID, movieID)
	if err != nil {
		h.Log.LogError(ctx, "ratings", "CreateRating", err)
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	ctx.Status(http.StatusOK)
}
