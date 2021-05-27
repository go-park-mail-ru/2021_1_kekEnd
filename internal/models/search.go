package models

// SearchResult структура результата поиска
type SearchResult struct {
	Movies []Movie `json:"movies"`
	Actors []Actor `json:"actors"`
	Users  []User  `json:"users"`
}
