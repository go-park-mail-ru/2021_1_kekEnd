package usecase

import (
	"errors"
	"testing"

	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/playlists/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestPlaylistsUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockPlaylistsRepository(ctrl)
	uc := NewPlaylistsUseCase(repo)

	t.Run("CreatePlaylist", func(t *testing.T) {
		repo.EXPECT().CreatePlaylist("let_robots_reign", "playlist", false).Return(nil)
		err := uc.CreatePlaylist("let_robots_reign", "playlist", false)
		assert.NoError(t, err)
	})

	t.Run("GetPlaylist", func(t *testing.T) {
		playlist := &models.Playlist{
			ID:       "1",
			Name:     "playlist",
			IsShared: false,
			Movies:   nil,
		}
		repo.EXPECT().GetPlaylist(1).Return(playlist, nil)
		gotPlaylist, err := uc.GetPlaylist(1)
		assert.NoError(t, err)
		assert.Equal(t, gotPlaylist, playlist)
	})

	t.Run("GetPlaylists", func(t *testing.T) {
		playlists := []models.Playlist{
			models.Playlist{
				ID:       "1",
				Name:     "playlist",
				IsShared: false,
				Movies:   nil,
			},
		}
		repo.EXPECT().GetPlaylists("user").Return(playlists, nil)
		gotPlaylists, err := uc.GetPlaylists("user")
		assert.NoError(t, err)
		assert.Equal(t, gotPlaylists, playlists)
	})

	t.Run("GetPlaylistsInfo", func(t *testing.T) {
		playlistsInfo := []models.PlaylistsInfo{
			{
				ID:      "1",
				Name:    "test",
				IsAdded: false,
			},
		}
		repo.EXPECT().GetPlaylistsInfo("user", 1).Return(playlistsInfo, nil)
		gotPlaylists, err := uc.GetPlaylistsInfo("user", 1)
		assert.NoError(t, err)
		assert.Equal(t, gotPlaylists, playlistsInfo)
	})

	t.Run("UpdatePlaylist", func(t *testing.T) {
		repo.EXPECT().CanUserUpdatePlaylist("user", 1).Return(nil)
		repo.EXPECT().UpdatePlaylist("user", 1, "test", true)
		err := uc.UpdatePlaylist("user", 1, "test", true)
		assert.NoError(t, err)
	})

	t.Run("UpdatePlaylist2", func(t *testing.T) {
		repo.EXPECT().CanUserUpdatePlaylist("user", 1).Return(nil)
		repo.EXPECT().UpdatePlaylist("user", 1, "test", false)
		repo.EXPECT().DeleteAllUserFromPlaylist("user", 1)
		err := uc.UpdatePlaylist("user", 1, "test", false)
		assert.NoError(t, err)
	})

	t.Run("DeletePlaylist", func(t *testing.T) {
		repo.EXPECT().CanUserUpdatePlaylist("user", 1).Return(nil)
		repo.EXPECT().DeletePlaylist(1).Return(nil)
		err := uc.DeletePlaylist("user", 1)
		assert.NoError(t, err)
	})

	t.Run("AddMovieToPlaylist", func(t *testing.T) {
		repo.EXPECT().CanUserUpdateMovieInPlaylist("user", 1).Return(nil)
		repo.EXPECT().AddMovieToPlaylist("user", 1, 1).Return(nil)
		err := uc.AddMovieToPlaylist("user", 1, 1)
		assert.NoError(t, err)
	})

	t.Run("DeleteMovieFromPlaylist", func(t *testing.T) {
		repo.EXPECT().CanUserUpdateMovieInPlaylist("user", 1).Return(nil)
		repo.EXPECT().DeleteMovieFromPlaylist("user", 1, 1).Return(nil)
		err := uc.DeleteMovieFromPlaylist("user", 1, 1)
		assert.NoError(t, err)
	})

	t.Run("AddUserToPlaylist", func(t *testing.T) {
		repo.EXPECT().CanUserUpdateUsersInPlaylist("user", 1).Return(nil)
		repo.EXPECT().AddUserToPlaylist("user", 1, "bob").Return(nil)
		err := uc.AddUserToPlaylist("user", 1, "bob")
		assert.NoError(t, err)
	})

	t.Run("DeleteUserFromPlaylist", func(t *testing.T) {
		repo.EXPECT().CanUserUpdateUsersInPlaylist("user", 1).Return(nil)
		repo.EXPECT().DeleteUserFromPlaylist("user", 1, "bob").Return(nil)
		err := uc.DeleteUserFromPlaylist("user", 1, "bob")
		assert.NoError(t, err)
	})

	t.Run("UpdatePlaylistError", func(t *testing.T) {
		repo.EXPECT().CanUserUpdatePlaylist("user", 1).Return(errors.New("error"))
		err := uc.UpdatePlaylist("user", 1, "test", true)
		assert.Error(t, err)
	})

	t.Run("UpdatePlaylistError", func(t *testing.T) {
		repo.EXPECT().CanUserUpdatePlaylist("user", 1).Return(nil)
		repo.EXPECT().UpdatePlaylist("user", 1, "test", false).Return(errors.New("error"))
		err := uc.UpdatePlaylist("user", 1, "test", false)
		assert.Error(t, err)
	})

	t.Run("UpdatePlaylistError", func(t *testing.T) {
		repo.EXPECT().CanUserUpdatePlaylist("user", 1).Return(nil)
		repo.EXPECT().UpdatePlaylist("user", 1, "test", false).Return(nil)
		repo.EXPECT().DeleteAllUserFromPlaylist("user", 1).Return(errors.New("error"))
		err := uc.UpdatePlaylist("user", 1, "test", false)
		assert.Error(t, err)
	})

	t.Run("DeletePlaylistError", func(t *testing.T) {
		repo.EXPECT().CanUserUpdatePlaylist("user", 1).Return(errors.New("error"))
		err := uc.DeletePlaylist("user", 1)
		assert.Error(t, err)
	})

	t.Run("AddMovieToPlaylistError", func(t *testing.T) {
		repo.EXPECT().CanUserUpdateMovieInPlaylist("user", 1).Return(errors.New("error"))
		err := uc.AddMovieToPlaylist("user", 1, 1)
		assert.Error(t, err)
	})

	t.Run("DeleteMovieFromPlaylistError", func(t *testing.T) {
		repo.EXPECT().CanUserUpdateMovieInPlaylist("user", 1).Return(errors.New("error"))
		err := uc.DeleteMovieFromPlaylist("user", 1, 1)
		assert.Error(t, err)
	})

	t.Run("AddUserToPlaylistError", func(t *testing.T) {
		repo.EXPECT().CanUserUpdateUsersInPlaylist("user", 1).Return(errors.New("error"))
		err := uc.AddUserToPlaylist("user", 1, "bob")
		assert.Error(t, err)
	})

	t.Run("DeleteUserFromPlaylistError", func(t *testing.T) {
		repo.EXPECT().CanUserUpdateUsersInPlaylist("user", 1).Return(errors.New("error"))
		err := uc.DeleteUserFromPlaylist("user", 1, "bob")
		assert.Error(t, err)
	})
}
