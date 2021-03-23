package models

type Rating struct {
	UserID  string `json:"user_id"`
	MovieID string `json:"movie_id"`
	Score   int    `json:"score"`
}
