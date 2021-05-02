package models

type Playlist struct {
	ID       string           `json:"id"`
	Name     string           `json:"name"`
	IsShared bool             `json:"isShared"`
	Movies   []MovieReference `json:"movies"`
}

type MovieInPlaylist struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	AddedBy string `json:"username"`
}

type PlaylistsInfo struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	IsAdded bool   `json:"isAdded"`
}
