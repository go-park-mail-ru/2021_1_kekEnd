package _const

import "time"

const (
	UserKey           = "user"
	ActorKey          = "actor"
	AuthStatusKey     = "auth_status"
	CookieExpires     = 240 * time.Hour
	CsrfExpires       = 10 * time.Minute
	Host              = "localhost"
	Port              = "8080"
	AvatarsPath       = "http://" + Host + ":" + Port + "/avatars/"
	DefaultAvatarPath = "http://" + Host + ":" + Port + "/avatars/default.jpeg"
	AvatarsFileDir    = "tmp/avatars/"
	//TODO
	RequestID = "RequestID"

	ReviewsPageSize         = 3
	SubsPageSize            = 20
	MoviesPageSize          = 15
	MoviesTop100Size        = 100
	MoviesNumberOnActorPage = 10
	PageDefault             = "1"
)

var AdminUsers = []string{
	"let_robots_reign",
	"IfuryI",
	"grillow",
	"polyanimal",
}
