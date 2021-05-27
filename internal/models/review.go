package models

import "time"

type ReviewType string

type Review struct {
	ID           string     `json:"id"`
	Title        string     `json:"title"`
	ReviewType   ReviewType `json:"review_type"`
	Content      string     `json:"content"`
	Author       string     `json:"author"`
	MovieID      string     `json:"movie_id"`
	CreationDate time.Time  `json:"creation_date,omitempty"`
}
