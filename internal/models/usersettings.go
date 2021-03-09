package models

type UserSettings struct {
	Username string `json:"username"`
	Email    string `json:"email"`

	MoviesWatched uint `json:"movies_watched"`
	ReviewsNumber uint `json:"reviews_number"`
}
