package playlists

import (
	"errors"

	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/playlists"
)

type PlaylistUseCase struct {
	playlistRepository playlists.Repository
}

func NewPlaylistUseCase(playlistRepo playlists.Repository) *PlaylistUseCase {
	return &PlaylistUseCase{
		playlistRepository: playlistRepo,
	}
}

func (playlistUC *PlaylistUseCase) CreatePlaylist(username string, playlistID int, isShared bool) error {
	_, err := playlistUC.GetPlaylistsInfo(playlistID)
	if err == nil {
		return errors.New("playlist already exists")
	}

	err = playlistUC.playlistRepository.CreatePlaylist(username, playlistID, isShared)
	if err != nil {
		return err
	}

	return err
}

func (playlistUC *PlaylistUseCase) GetPlaylistsInfo(playlistID int) (*models.Playlist, error) {
	return playlistUC.playlistRepository.GetPlaylistsInfo(playlistID)
}

func (playlistUC *PlaylistUseCase) UpdatePlaylist(username string, playlistID int, isShared bool) (*models.Playlist, error) {
	return playlistUC.playlistRepository.UpdatePlaylist(username, playlistID, isShared)
}

func (playlistUC *PlaylistUseCase) DeletePlaylist(username string, playlistID int) error {
	return playlistUC.playlistRepository.DeletePlaylist(username, playlistID)
}

func (playlistUC *PlaylistUseCase) AddToPlaylist(username string, playlistData *models.Playlist) error {
	// oldReview, err := reviewsUC.GetUserReviewForMovie(user.Username, review.MovieID)
	// if err != nil {
	// 	return errors.New("review doesn't exist")
	// }
	// review.ID = oldReview.ID
	// review.Author = user.Username
	return playlistUC.playlistRepository.AddToPlaylist(username, playlistData)
}

func (playlistUC *PlaylistUseCase) DeleteFromPlaylist(username string, playlistID int, movieID int) error {
	return playlistUC.playlistRepository.DeleteFromPlaylist(username, playlistID, movieID)
}
