package playlists

import "github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"

//go:generate mockgen -destination=mocks/usecase.go -package=mocks . UseCase
type UseCase interface {
	CreatePlaylist(username string, playlistName string, isShared bool) error
	GetPlaylistsInfo(username string, movieID int) ([]*models.PlaylistsInfo, error)
	GetPlaylists(username string) ([]*models.Playlist, error)
	UpdatePlaylist(username string, playlistID int, isShared bool) (*models.Playlist, error)
	DeletePlaylist(username string, playlistID int) error

	AddToPlaylist(username string, playlistData *models.Playlist) error
	DeleteFromPlaylist(username string, playlistID int, movieID int) error
}
