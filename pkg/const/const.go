package _const

import "time"

const (
	UserKey           = "user"
	ActorKey          = "actor"
	CookieExpires     = 240 * time.Hour
	Host              = "localhost"
	Port              = "8080"
	AvatarsPath       = "http://" + Host + ":" + Port + "/avatars/"
	DefaultAvatarPath = "http://" + Host + ":" + Port + "/avatars/default.jpeg"
	AvatarsFileDir    = "tmp/avatars/"

	ReviewsPageSize         = 3
	MoviesPageSize          = 15
	MoviesTop100Size        = 100
	MoviesNumberOnActorPage = 10
	PageDefault             = "1"
)
