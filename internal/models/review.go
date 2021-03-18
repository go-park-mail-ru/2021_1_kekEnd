package models

type reviewType string

const (
	negative reviewType = "negative"
	neutral             = "neutral"
	positive            = "positive"
)

type Review struct {
	ID         string     `json:"id"`
	Title      string     `json:"title"`
	ReviewType reviewType `json:"review_type"`
	Content    string     `json:"content"`
	Author     string     `json:"author"`
	MovieID    string     `json:"movie_id"`
}
