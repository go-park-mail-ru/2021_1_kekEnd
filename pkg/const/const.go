package _const

import "time"

const (
	UserKey           = "user"
	CookieExpires     = 240 * time.Hour
	Host              = "localhost"
	Port              = "8080"
	AvatarsPath       = "http://" + Host + ":" + Port + "/avatars/"
	DefaultAvatarPath = "http://" + Host + ":" + Port + "/avatars/default.jpeg"
	AvatarsFileDir    = "tmp/avatars/"

	// TODO: should be 10, but for testing purposes is 3
	ReviewsPageSize = 3

	PageDefault = "1"
)
