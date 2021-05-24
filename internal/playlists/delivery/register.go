package http

import (
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/logger"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/middleware"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/playlists"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/users"
)

func RegisterHttpEndpoints(router *gin.RouterGroup, playlistsUC playlists.UseCase, usersUC users.UseCase,
	authMiddleware middleware.Auth, Log *logger.Logger) {
	handler := NewHandler(playlistsUC, usersUC, Log)

	router.POST("/playlists", authMiddleware.CheckAuth(true), handler.CreatePlaylist)
	router.GET("/playlist/:playlist_id", handler.GetPlaylist)
	router.GET("/playlists/movies/:movie_id", authMiddleware.CheckAuth(true), handler.GetPlaylistsInfo)
	router.GET("/playlists/users/:username", handler.GetPlaylists)
	router.PUT("/playlists", authMiddleware.CheckAuth(true), handler.EditPlaylist)
	router.DELETE("/playlists/:playlist_id", authMiddleware.CheckAuth(true), handler.DeletePlaylist)

	router.POST("/playlists/:playlist_id/movie", authMiddleware.CheckAuth(true), handler.AddMovieToPlaylist)
	router.DELETE("/playlists/:playlist_id/movie", authMiddleware.CheckAuth(true), handler.DeleteMovieFromPlaylist)

	router.POST("/playlists/:playlist_id/user", authMiddleware.CheckAuth(true), handler.AddUserToPlaylist)
	router.DELETE("/playlists/:playlist_id/user", authMiddleware.CheckAuth(true), handler.DeleteUserFromPlaylist)
}
