package models

const (
	NewReview = iota
	NewRating
)

type Notification struct {
	Type    int
	Message string
}
