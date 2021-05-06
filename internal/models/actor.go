package models

type Actor struct {
	ID           string           `json:"id"`
	Name         string           `json:"name"`
	Biography    string           `json:"biography"`
	BirthDate    string           `json:"birthdate"`
	Origin       string           `json:"origin"`
	Profession   string           `json:"profession"`
	MoviesCount  int              `json:"movies_count"`
	MoviesRating int              `json:"movies_rating"`
	Movies       []MovieReference `json:"movies"`
	Avatar       string           `json:"avatar"`
	IsLiked      bool             `json:"is_liked"`
}

type MovieReference struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Rating float64 `json:"rating"`
}
