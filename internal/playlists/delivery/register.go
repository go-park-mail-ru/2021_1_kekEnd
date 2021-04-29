package playlists

import (
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/logger"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/middleware"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/playlists"
)

func RegisterHttpEndpoints(router *gin.Engine, playlistsUC playlists.UseCase, authMiddleware middleware.Auth, Log *logger.Logger) {
	handler := NewHandler(playlistsUC, Log)

	router.POST("/playlists", authMiddleware.CheckAuth(), handler.CreatePlaylist)
	router.POST("/playlists/:playlist_id", authMiddleware.CheckAuth(), handler.AddToPlaylist)
	router.GET("/playlists/:playlist_id", handler.GetPlaylist)
	router.PUT("/playlists", authMiddleware.CheckAuth(), handler.EditPlaylist)
	router.DELETE("/playlists/:playlist_id", authMiddleware.CheckAuth(), handler.DeletePlaylist)
	router.DELETE("/playlists/:playlist_id", authMiddleware.CheckAuth(), handler.DeleteFromPlaylist)
}
