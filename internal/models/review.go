package models

type reviewType string

const (
	negative reviewType = "negative"
	neutral = "neutral"
	positive = "positive"
)

type Review struct {
	Title      string     `json:"title"`
	ReviewType reviewType `json:"review_type"`
	Content    string     `json:"content"`
}
