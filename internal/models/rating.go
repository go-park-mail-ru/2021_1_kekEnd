package models

type Rating struct {
	UserID  string `json:"username"`
	MovieID string `json:"movie_id"`
	Score   int    `json:"score"`
}
