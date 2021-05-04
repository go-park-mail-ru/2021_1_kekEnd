package models

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Avatar   string `json:"avatar"`

	MoviesWatched *uint `json:"movies_watched"`
	ReviewsNumber *uint `json:"reviews_number"`

	FavoriteActors []Actor `json:"favorite_actors"`
}

type UserNoPassword struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`

	MoviesWatched *uint `json:"movies_watched"`
	ReviewsNumber *uint `json:"reviews_number"`

	FavoriteActors []Actor `json:"favorite_actors"`
}

func FromUser(user User) UserNoPassword {
	return UserNoPassword{
		Username:      user.Username,
		Email:         user.Email,
		Avatar:        user.Avatar,
		MoviesWatched: user.MoviesWatched,
		ReviewsNumber: user.ReviewsNumber,
		FavoriteActors: user.FavoriteActors,
	}
}
