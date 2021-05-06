package usecase

import (
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/playlists/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
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
}