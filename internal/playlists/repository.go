package playlists

import "github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"

// Repository go:generate mockgen -destination=mocks/repository.go -package=mocks . Repository
type Repository interface {
	CreatePlaylist(username string, playlistName string, isShared bool) error

	GetPlaylist(playlistID int) (*models.Playlist, error)
	GetPlaylists(username string) ([]models.Playlist, error)
	GetPlaylistsInfo(username string, movieID int) ([]models.PlaylistsInfo, error)

	CanUserUpdatePlaylist(username string, playlistID int) error
	UpdatePlaylist(username string, playlistID int, playlistName string, isShared bool) error

	DeletePlaylist(playlistID int) error

	CanUserUpdateMovieInPlaylist(username string, playlistID int) error
	AddMovieToPlaylist(username string, playlistID int, movieID int) error
	DeleteMovieFromPlaylist(username string, playlistID int, movieID int) error

	CanUserUpdateUsersInPlaylist(username string, playlistID int) error
	AddUserToPlaylist(username string, playlistID int, usernameToAdd string) error
	DeleteUserFromPlaylist(username string, playlistID int, usernameToDelete string) error
	DeleteAllUserFromPlaylist(username string, playlistID int) error
}
