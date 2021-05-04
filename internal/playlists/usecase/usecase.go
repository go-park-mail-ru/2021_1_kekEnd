package usecase

import (
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/playlists"
)

type PlaylistUseCase struct {
	playlistRepository playlists.PlaylistsRepository
}

func NewPlaylistsUseCase(playlistRepo playlists.PlaylistsRepository) *PlaylistUseCase {
	return &PlaylistUseCase{
		playlistRepository: playlistRepo,
	}
}

func (playlistUC *PlaylistUseCase) CreatePlaylist(username string, playlistName string, isShared bool) error {
	return playlistUC.playlistRepository.CreatePlaylist(username, playlistName, isShared)
}

func (playlistUC *PlaylistUseCase) GetPlaylist(playlistID int) (*models.Playlist, error) {
	return playlistUC.playlistRepository.GetPlaylist(playlistID)
}

func (playlistUC *PlaylistUseCase) GetPlaylists(username string) ([]*models.Playlist, error) {
	return playlistUC.playlistRepository.GetPlaylists(username)
}

func (playlistUC *PlaylistUseCase) GetPlaylistsInfo(username string, movieID int) ([]*models.PlaylistsInfo, error) {
	return playlistUC.playlistRepository.GetPlaylistsInfo(username, movieID)
}

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

func (playlistUC *PlaylistUseCase) DeletePlaylist(username string, playlistID int) error {
	err := playlistUC.playlistRepository.CanUserUpdatePlaylist(username, playlistID)
	if err != nil {
		return err
	}

	return playlistUC.playlistRepository.DeletePlaylist(playlistID)
}

func (playlistUC *PlaylistUseCase) AddMovieToPlaylist(username string, playlistID int, movieID int) error {
	err := playlistUC.playlistRepository.CanUserUpdateMovieInPlaylist(username, playlistID)
	if err != nil {
		return err
	}

	return playlistUC.playlistRepository.AddMovieToPlaylist(username, playlistID, movieID)
}

func (playlistUC *PlaylistUseCase) DeleteMovieFromPlaylist(username string, playlistID int, movieID int) error {
	err := playlistUC.playlistRepository.CanUserUpdateMovieInPlaylist(username, playlistID)
	if err != nil {
		return err
	}

	return playlistUC.playlistRepository.DeleteMovieFromPlaylist(username, playlistID, movieID)
}

func (playlistUC *PlaylistUseCase) AddUserToPlaylist(username string, playlistID int, usernameToAdd string) error {
	err := playlistUC.playlistRepository.CanUserUpdateUsersInPlaylist(username, playlistID)
	if err != nil {
		return err
	}

	return playlistUC.playlistRepository.AddUserToPlaylist(username, playlistID, usernameToAdd)
}

func (playlistUC *PlaylistUseCase) DeleteUserFromPlaylist(username string, playlistID int, usernameToDelete string) error {
	err := playlistUC.playlistRepository.CanUserUpdateUsersInPlaylist(username, playlistID)
	if err != nil {
		return err
	}

	return playlistUC.playlistRepository.DeleteUserFromPlaylist(username, playlistID, usernameToDelete)
}
