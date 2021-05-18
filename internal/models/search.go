package models

type SearchResult struct {
	Movies []Movie `json:"movies"`
	Actors []Actor `json:"actors"`
	Users  []User  `json:"users"`
}
