package playlists

import "github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"

//go:generate mockgen -destination=mocks/repository.go -package=mocks . Repository
type Repository interface {
	CreatePlaylist(username string, playlistID int, isShared bool) error
	GetPlaylistsInfo(playlistID int) (*models.Playlist, error)
	UpdatePlaylist(username string, playlistID int, isShared bool) (*models.Playlist, error)
	DeletePlaylist(username string, playlistID int) error

	AddToPlaylist(username string, playlistData *models.Playlist) error
	DeleteFromPlaylist(username string, playlistID int, movieID int) error
}
