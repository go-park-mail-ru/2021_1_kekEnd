package models

type Playlist struct {
	ID     string           `json:"id"`
	Name   string           `json:"name"`
	Movies []MovieReference `json:"movies"`
}

type MovieReference struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	AddedBy string `json:"username"`
}
