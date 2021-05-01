package playlists

import (
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/logger"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/middleware"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/playlists"
)

func RegisterHttpEndpoints(router *gin.Engine, playlistsUC playlists.UseCase, authMiddleware middleware.Auth, Log *logger.Logger) {
	handler := NewHandler(playlistsUC, Log)

	router.POST("/users/playlists", authMiddleware.CheckAuth(), handler.CreatePlaylist)
	router.GET("/users/playlists/info", authMiddleware.CheckAuth(), handler.GetPlaylistsInfo)
	router.GET("/users/playlists", authMiddleware.CheckAuth(), handler.GetPlaylists)
	router.PUT("/users/playlists", authMiddleware.CheckAuth(), handler.EditPlaylist)
	router.DELETE("/users/playlists/:playlist_id", authMiddleware.CheckAuth(), handler.DeletePlaylist)

	router.POST("/users/playlists/:playlist_id", authMiddleware.CheckAuth(), handler.AddToPlaylist)
	router.DELETE("/users/playlists", authMiddleware.CheckAuth(), handler.DeleteFromPlaylist)
}
