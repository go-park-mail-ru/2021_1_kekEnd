package http

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

type PlaylistMovie struct {
	PlaylistID string `json:"playlistID"`
	MovieID    string `json:"movieID"`
}

type PlaylistUser struct {
	PlaylistID string `json:"playlistID"`
	Username   string `json:"username"`
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

	fmt.Println(playlistData)

	// user, ok := ctx.Get(_const.UserKey)
	// if !ok {
	// 	err := fmt.Errorf("%s", "Failed to retrieve user from context")
	// 	h.Log.LogWarning(ctx, "playlists", "CreatePlaylist", err.Error())
	// 	ctx.AbortWithStatus(http.StatusBadRequest) // 400
	// 	return
	// }

	// userModel, ok := user.(models.User)
	// if !ok {
	// 	err := fmt.Errorf("%s", "Failed to cast user to model")
	// 	h.Log.LogError(ctx, "playlists", "CreatePlaylist", err)
	// 	ctx.AbortWithStatus(http.StatusInternalServerError) // 500
	// }

	err = h.useCase.CreatePlaylist("user1", playlistData.Name, playlistData.IsShared)
	if err != nil {
		h.Log.LogError(ctx, "playlists", "CreatePlaylist", err)
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	ctx.Status(http.StatusCreated)
}

func (h *Handler) GetPlaylistsInfo(ctx *gin.Context) {
	movieIDStr := ctx.Param("movie_id")
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

	movieID, err := strconv.Atoi(movieIDStr)
	if err != nil {
		err := fmt.Errorf("%s", "Failed to cast rating value to number")
		h.Log.LogWarning(ctx, "ratings", "CreateRating", err.Error())
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	playlistInfo, err := h.useCase.GetPlaylistsInfo(userModel.Username, movieID)
	if err != nil {
		h.Log.LogWarning(ctx, "ratings", "GetRating", err.Error())
		ctx.AbortWithStatus(http.StatusNotFound) // 404
		return
	}

	ctx.JSON(http.StatusOK, playlistInfo)
}

func (h *Handler) GetPlaylists(ctx *gin.Context) {
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

	playlists, err := h.useCase.GetPlaylists(userModel.Username)
	if err != nil {
		h.Log.LogWarning(ctx, "ratings", "GetRating", err.Error())
		ctx.AbortWithStatus(http.StatusNotFound) // 404
		return
	}

	ctx.JSON(http.StatusOK, playlists)
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

	playlist, err := h.useCase.UpdatePlaylist(userModel.Username, playlistID, playlistData.Name, playlistData.IsShared)
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

func (h *Handler) AddMovieToPlaylist(ctx *gin.Context) {
	playlistMovieData := new(PlaylistMovie)
	err := ctx.BindJSON(playlistMovieData)
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

	playlistID, err := strconv.Atoi(playlistMovieData.PlaylistID)
	if err != nil {
		err := fmt.Errorf("%s", "Failed to cast rating value to number")
		h.Log.LogWarning(ctx, "ratings", "CreateRating", err.Error())
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	movieID, err := strconv.Atoi(playlistMovieData.MovieID)
	if err != nil {
		err := fmt.Errorf("%s", "Failed to cast rating value to number")
		h.Log.LogWarning(ctx, "ratings", "CreateRating", err.Error())
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	err = h.useCase.AddMovieToPlaylist(userModel.Username, playlistID, movieID)
	if err != nil {
		h.Log.LogError(ctx, "ratings", "CreateRating", err)
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	ctx.Status(http.StatusCreated)
}

func (h *Handler) DeleteMovieFromPlaylist(ctx *gin.Context) {
	playlistMovieData := new(PlaylistMovie)
	err := ctx.BindJSON(playlistMovieData)
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

	playlistID, err := strconv.Atoi(playlistMovieData.PlaylistID)
	if err != nil {
		err := fmt.Errorf("%s", "Failed to cast rating value to number")
		h.Log.LogWarning(ctx, "ratings", "CreateRating", err.Error())
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	movieID, err := strconv.Atoi(playlistMovieData.MovieID)
	if err != nil {
		err := fmt.Errorf("%s", "Failed to cast rating value to number")
		h.Log.LogWarning(ctx, "ratings", "CreateRating", err.Error())
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	err = h.useCase.DeleteMovieFromPlaylist(userModel.Username, playlistID, movieID)
	if err != nil {
		h.Log.LogError(ctx, "ratings", "CreateRating", err)
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	ctx.Status(http.StatusOK)
}

func (h *Handler) AddUserToPlaylist(ctx *gin.Context) {
	playlistUserData := new(PlaylistUser)
	err := ctx.BindJSON(playlistUserData)
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

	playlistID, err := strconv.Atoi(playlistUserData.PlaylistID)
	if err != nil {
		err := fmt.Errorf("%s", "Failed to cast rating value to number")
		h.Log.LogWarning(ctx, "ratings", "CreateRating", err.Error())
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	err = h.useCase.AddUserToPlaylist(userModel.Username, playlistID, playlistUserData.Username)
	if err != nil {
		h.Log.LogError(ctx, "ratings", "CreateRating", err)
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	ctx.Status(http.StatusCreated)
}

func (h *Handler) DeleteUserFromPlaylist(ctx *gin.Context) {
	playlistUserData := new(PlaylistUser)
	err := ctx.BindJSON(playlistUserData)
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

	playlistID, err := strconv.Atoi(playlistUserData.PlaylistID)
	if err != nil {
		err := fmt.Errorf("%s", "Failed to cast rating value to number")
		h.Log.LogWarning(ctx, "ratings", "CreateRating", err.Error())
		ctx.AbortWithStatus(http.StatusBadRequest) // 400
		return
	}

	err = h.useCase.DeleteUserFromPlaylist(userModel.Username, playlistID, playlistUserData.Username)
	if err != nil {
		h.Log.LogError(ctx, "ratings", "CreateRating", err)
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	ctx.Status(http.StatusOK)
}
