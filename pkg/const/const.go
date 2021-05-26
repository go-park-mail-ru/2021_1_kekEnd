package _const

import "time"

const (
	UserKey           = "user"
	ActorKey          = "actor"
	AuthStatusKey     = "auth_status"
	CookieExpires     = 240 * time.Hour
	CsrfExpires       = 10 * time.Minute
	Host              = "cinemedia.ru"
	Port              = "8085"
	AuthPort          = "8081"
	FileServerPort    = "8082"
	RedisPort         = "6379"
	AvatarsPath       = "https://" + Host + "/avatars/"
	DefaultAvatarPath = "https://" + Host + "/avatars/default.jpeg"
	AvatarsFileDir    = "tmp/avatars/"
	PostersFileDir    = "tmp/posters/"
	BannersFileDir    = "tmp/banners/"
	ActorsFileDir     = "tmp/actors/"
	//TODO
	RequestID = "RequestID"

	ReviewsPageSize         = 3
	SubsPageSize            = 20
	MoviesPageSize          = 15
	MoviesTop100Size        = 100
	MoviesNumberOnActorPage = 10
	FeedItemsLimit          = 20
	SimilarMoviesLimit      = 10
	PageDefault             = "1"
)

var AdminUsers = []string{
	"let_robots_reign",
	"IfuryI",
	"grillow",
	"polyanimal",
}
