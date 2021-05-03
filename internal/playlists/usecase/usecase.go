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
	// _, err := playlistUC.GetPlaylistsInfo(username, )
	// if err == nil {
	// 	return errors.New("playlist already exists")
	// }

	err := playlistUC.playlistRepository.CreatePlaylist(username, playlistName, isShared)
	if err != nil {
		return err
	}

	return err
}

func (playlistUC *PlaylistUseCase) GetPlaylistsInfo(username string, movieID int) ([]*models.PlaylistsInfo, error) {
	return playlistUC.playlistRepository.GetPlaylistsInfo(username, movieID)
}

func (playlistUC *PlaylistUseCase) GetPlaylists(username string) ([]*models.Playlist, error) {
	return playlistUC.playlistRepository.GetPlaylists(username)
}

func (playlistUC *PlaylistUseCase) UpdatePlaylist(username string, playlistID int, playlistName string, isShared bool) error {
	return playlistUC.playlistRepository.UpdatePlaylist(username, playlistID, playlistName, isShared)
}

func (playlistUC *PlaylistUseCase) DeletePlaylist(username string, playlistID int) error {
	return playlistUC.playlistRepository.DeletePlaylist(playlistID)
}

func (playlistUC *PlaylistUseCase) AddMovieToPlaylist(username string, playlistID int, movieID int) error {
	// oldReview, err := reviewsUC.GetUserReviewForMovie(user.Username, review.MovieID)
	// if err != nil {
	// 	return errors.New("review doesn't exist")
	// }
	// review.ID = oldReview.ID
	// review.Author = user.Username
	return playlistUC.playlistRepository.AddMovieToPlaylist(username, playlistID, movieID)
}

func (playlistUC *PlaylistUseCase) DeleteMovieFromPlaylist(username string, playlistID int, movieID int) error {
	return playlistUC.playlistRepository.DeleteMovieFromPlaylist(username, playlistID, movieID)
}

func (playlistUC *PlaylistUseCase) AddUserToPlaylist(username string, playlistID int, usernameToAdd string) error {
	// oldReview, err := reviewsUC.GetUserReviewForMovie(user.Username, review.MovieID)
	// if err != nil {
	// 	return errors.New("review doesn't exist")
	// }
	// review.ID = oldReview.ID
	// review.Author = user.Username
	return playlistUC.playlistRepository.AddUserToPlaylist(username, playlistID, usernameToAdd)
}

func (playlistUC *PlaylistUseCase) DeleteUserFromPlaylist(username string, playlistID int, usernameToDelete string) error {
	return playlistUC.playlistRepository.DeleteUserFromPlaylist(username, playlistID, usernameToDelete)
}
