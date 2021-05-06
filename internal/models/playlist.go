package models

type Playlist struct {
	ID       string            `json:"id"`
	Name     string            `json:"playlist_name"`
	IsShared bool              `json:"is_shared"`
	Movies   []MovieInPlaylist `json:"movies,omitempty"`
}

type MovieInPlaylist struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	AddedBy string `json:"username,omitempty"`
}

type PlaylistsInfo struct {
	ID      string `json:"id"`
	Name    string `json:"playlist_name"`
	IsAdded bool   `json:"is_added"`
}
