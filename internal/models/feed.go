package models

import "time"

// ReviewFeedItem структура новостной рецензии
type ReviewFeedItem struct {
	Username string    `json:"username"`
	Avatar   string    `json:"avatar"`
	ItemType string    `json:"item_type"`
	Review   Review    `json:"review"`
	Date     time.Time `json:"date"`
}

// RatingFeedItem структура новостной оценки
type RatingFeedItem struct {
	Username   string    `json:"username"`
	Avatar     string    `json:"avatar"`
	ItemType   string    `json:"item_type"`
	MovieTitle string    `json:"movie_title"`
	Rating     Rating    `json:"rating"`
	Date       time.Time `json:"date"`
}

// Feed структура новостей
type Feed struct {
	Ratings []RatingFeedItem `json:"recent_ratings"`
	Reviews []ReviewFeedItem `json:"recent_reviews"`
}
