package models

type Playlist struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	IsShared bool              `json:"isShared"`
	Movies   []MovieInPlaylist `json:"movies,omitempty"`
}

type MovieInPlaylist struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	AddedBy string `json:"username,omitempty"`
}

type PlaylistsInfo struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	IsAdded bool   `json:"isAdded"`
}
