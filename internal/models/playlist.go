package models

// Playlist структура плейлиста
type Playlist struct {
	ID       string            `json:"id"`
	Name     string            `json:"playlist_name"`
	IsShared bool              `json:"is_shared"`
	Movies   []MovieInPlaylist `json:"movies,omitempty"`
}

// MovieInPlaylist структура фильма в плейлисте
type MovieInPlaylist struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	AddedBy string `json:"username,omitempty"`
}

// PlaylistsInfo структура информации о плейлисте
type PlaylistsInfo struct {
	ID      string `json:"id"`
	Name    string `json:"playlist_name"`
	IsAdded bool   `json:"is_added"`
}
