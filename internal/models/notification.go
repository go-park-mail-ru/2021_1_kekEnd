package models

import "time"

type ReviewFeedItem struct {
	Username string    `json:"username"`
	Avatar   string    `json:"avatar"`
	ItemType string    `json:"item_type"`
	Review   Review    `json:"review"`
	Date     time.Time `json:"date"`
}

type RatingFeedItem struct {
	Username string    `json:"username"`
	Avatar   string    `json:"avatar"`
	ItemType string    `json:"item_type"`
	Rating   Rating    `json:"rating"`
	Date     time.Time `json:"date"`
}
