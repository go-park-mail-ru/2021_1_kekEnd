package models

type User struct {
	ID       int
	Username string
	Email    string
	Password string

	FirstName string
	LastName string
	MoviesWatched uint
	ReviewsNumber uint
}
