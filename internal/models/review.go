package models

type ReviewType string

const (
	negative ReviewType = "negative"
	neutral             = "neutral"
	positive            = "positive"
)

type Review struct {
	ID         string     `json:"id"`
	Title      string     `json:"title"`
	ReviewType ReviewType `json:"review_type"`
	Content    string     `json:"content"`
	Author     string     `json:"author"`
	MovieID    string     `json:"movie_id"`
}
