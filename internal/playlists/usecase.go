package playlists

import "github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"

//go:generate mockgen -destination=mocks/usecase.go -package=mocks . UseCase
type UseCase interface {
	CreatePlaylist(userID string) error
	GetPlaylist(userID string, playlistID string) (models.Playlist, error)
	UpdatePlaylist(userID string, playlistID string) error
	DeletePlaylist(userID string, playlistID string) error

	AddToPlaylist(userID string, playlistID string) error
	DeleteFromPlaylist(userID string, playlistID string) error
}
