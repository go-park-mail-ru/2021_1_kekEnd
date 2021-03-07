package models

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`

	MoviesWatched uint `json:"movies_watched"`
	ReviewsNumber uint `json:"reviews_number"`
}
