package http

import (
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/logger"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/middleware"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/playlists"
)

func RegisterHttpEndpoints(router *gin.Engine, playlistsUC playlists.UseCase, authMiddleware middleware.Auth, Log *logger.Logger) {
	handler := NewHandler(playlistsUC, Log)

	router.POST("/playlists", handler.CreatePlaylist)
	router.GET("/playlists/:movie_id", authMiddleware.CheckAuth(), handler.GetPlaylistsInfo)
	router.GET("/playlists", authMiddleware.CheckAuth(), handler.GetPlaylists)
	router.PUT("/playlists", authMiddleware.CheckAuth(), handler.EditPlaylist)
	router.DELETE("/playlists/:playlist_id", authMiddleware.CheckAuth(), handler.DeletePlaylist)

	router.POST("/playlists1/movie", authMiddleware.CheckAuth(), handler.AddMovieToPlaylist)
	router.DELETE("/playlists1/movie", authMiddleware.CheckAuth(), handler.DeleteMovieFromPlaylist)

	router.POST("/playlists1/user", authMiddleware.CheckAuth(), handler.AddUserToPlaylist)
	router.DELETE("/playlists1/user", authMiddleware.CheckAuth(), handler.DeleteUserFromPlaylist)
}
