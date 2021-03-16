package models

type reviewType int

const (
	negative reviewType = iota
	neutral
	positive
)

type Review struct {
	Title      string     `json:"title"`
	ReviewType reviewType `json:"review_type"`
	Content    string     `json:"content"`
}
