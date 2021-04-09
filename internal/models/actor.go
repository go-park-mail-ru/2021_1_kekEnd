package models

type Actor struct {
	ID           string   `json:"id"`
	FirstName    string   `json:"first_name"`
	LastName     string   `json:"last_name"`
	Career       []string `json:"career"`
	Genres       []string `json:"genres"`
	PlaceOfBirth []string `json:"place_of_birth"`
	Birthday     string   `json:"birthday"`
	Height       string   `json:"height"`
	Spouse       string   `json:"spouse"`
	Photo        string   `json:"avatar"`
	BestMovies   []string `json:"best_movies"`
	BestSeries   []string `json:"best_series"`
	TotalMovies  uint     `json:"total_movies"`
}
