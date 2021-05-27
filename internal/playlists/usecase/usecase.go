package usecase

import (
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/playlists"
)

// PlaylistUseCase структура usecase плейлиста
type PlaylistUseCase struct {
	playlistRepository playlists.Repository
}

// NewPlaylistsUseCase инициализация usecase плейлиста
func NewPlaylistsUseCase(playlistRepo playlists.Repository) *PlaylistUseCase {
	return &PlaylistUseCase{
		playlistRepository: playlistRepo,
	}
}

// CreatePlaylist создание плейлиста
func (playlistUC *PlaylistUseCase) CreatePlaylist(username string, playlistName string, isShared bool) error {
	return playlistUC.playlistRepository.CreatePlaylist(username, playlistName, isShared)
}

// GetPlaylist получение плейлиста
func (playlistUC *PlaylistUseCase) GetPlaylist(playlistID int) (*models.Playlist, error) {
	return playlistUC.playlistRepository.GetPlaylist(playlistID)
}

// GetPlaylists получить все плейлисты
func (playlistUC *PlaylistUseCase) GetPlaylists(username string) ([]models.Playlist, error) {
	return playlistUC.playlistRepository.GetPlaylists(username)
}

// GetPlaylistsInfo получение информации о плейлисте
func (playlistUC *PlaylistUseCase) GetPlaylistsInfo(username string, movieID int) ([]models.PlaylistsInfo, error) {
	return playlistUC.playlistRepository.GetPlaylistsInfo(username, movieID)
}

// UpdatePlaylist изменить плейлист
func (playlistUC *PlaylistUseCase) UpdatePlaylist(username string, playlistID int, playlistName string, isShared bool) error {
	err := playlistUC.playlistRepository.CanUserUpdatePlaylist(username, playlistID)
	if err != nil {
		return err
	}

	err = playlistUC.playlistRepository.UpdatePlaylist(username, playlistID, playlistName, isShared)
	if err != nil {
		return err
	}

	if !isShared {
		err = playlistUC.playlistRepository.DeleteAllUserFromPlaylist(username, playlistID)
		if err != nil {
			return err
		}
	}

	return nil
}

// DeletePlaylist удалить плейлист
func (playlistUC *PlaylistUseCase) DeletePlaylist(username string, playlistID int) error {
	err := playlistUC.playlistRepository.CanUserUpdatePlaylist(username, playlistID)
	if err != nil {
		return err
	}

	return playlistUC.playlistRepository.DeletePlaylist(playlistID)
}

// AddMovieToPlaylist добавить фильм в плейлист
func (playlistUC *PlaylistUseCase) AddMovieToPlaylist(username string, playlistID int, movieID int) error {
	err := playlistUC.playlistRepository.CanUserUpdateMovieInPlaylist(username, playlistID)
	if err != nil {
		return err
	}

	return playlistUC.playlistRepository.AddMovieToPlaylist(username, playlistID, movieID)
}

// DeleteMovieFromPlaylist удалить фильм из плейлиста
func (playlistUC *PlaylistUseCase) DeleteMovieFromPlaylist(username string, playlistID int, movieID int) error {
	err := playlistUC.playlistRepository.CanUserUpdateMovieInPlaylist(username, playlistID)
	if err != nil {
		return err
	}

	return playlistUC.playlistRepository.DeleteMovieFromPlaylist(username, playlistID, movieID)
}

// AddUserToPlaylist добавить юзера в плейлист
func (playlistUC *PlaylistUseCase) AddUserToPlaylist(username string, playlistID int, usernameToAdd string) error {
	err := playlistUC.playlistRepository.CanUserUpdateUsersInPlaylist(username, playlistID)
	if err != nil {
		return err
	}

	return playlistUC.playlistRepository.AddUserToPlaylist(username, playlistID, usernameToAdd)
}

// DeleteUserFromPlaylist удалить юзера из плейлиста
func (playlistUC *PlaylistUseCase) DeleteUserFromPlaylist(username string, playlistID int, usernameToDelete string) error {
	err := playlistUC.playlistRepository.CanUserUpdateUsersInPlaylist(username, playlistID)
	if err != nil {
		return err
	}

	return playlistUC.playlistRepository.DeleteUserFromPlaylist(username, playlistID, usernameToDelete)
}
