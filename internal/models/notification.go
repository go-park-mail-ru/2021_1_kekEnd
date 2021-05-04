package models

import "time"

type Notification struct {
	Title string    `json:"title"`
	User  string    `json:"user"`
	Text  string    `json:"text"`
	Date  time.Time `json:"date"`
}
