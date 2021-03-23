package server

import (
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/middleware"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/movies"
	moviesHttp "github.com/go-park-mail-ru/2021_1_kekEnd/internal/movies/delivery/http"
	moviesLocalStorage "github.com/go-park-mail-ru/2021_1_kekEnd/internal/movies/repository/localstorage"
	moviesUseCase "github.com/go-park-mail-ru/2021_1_kekEnd/internal/movies/usecase"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/ratings"
	ratingsHttp "github.com/go-park-mail-ru/2021_1_kekEnd/internal/ratings/delivery"
	ratingsLocalStorage "github.com/go-park-mail-ru/2021_1_kekEnd/internal/ratings/repository/localstorage"
	ratingsUseCase "github.com/go-park-mail-ru/2021_1_kekEnd/internal/ratings/usecase"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/sessions"
	sessionsDelivery "github.com/go-park-mail-ru/2021_1_kekEnd/internal/sessions/delivery"
	sessionsRepository "github.com/go-park-mail-ru/2021_1_kekEnd/internal/sessions/repository"
	sessionsUseCase "github.com/go-park-mail-ru/2021_1_kekEnd/internal/sessions/usecase"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/users"
	usersHttp "github.com/go-park-mail-ru/2021_1_kekEnd/internal/users/delivery/http"
	usersLocalStorage "github.com/go-park-mail-ru/2021_1_kekEnd/internal/users/repository/localstorage"
	usersUseCase "github.com/go-park-mail-ru/2021_1_kekEnd/internal/users/usecase"
	_const "github.com/go-park-mail-ru/2021_1_kekEnd/pkg/const"
	"github.com/go-redis/redis/v8"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type App struct {
	server         *http.Server
	usersUC        users.UseCase
	moviesUC       movies.UseCase
	ratingsUC      ratings.UseCase
	sessions       sessions.Delivery
	authMiddleware middleware.Auth
}

func NewApp() *App {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	p, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Failed to create redis client", p, err)
	}

	sessionsRepo := sessionsRepository.NewRedisRepository(rdb)
	sessionsUC := sessionsUseCase.NewUseCase(sessionsRepo)
	sessionsDL := sessionsDelivery.NewDelivery(sessionsUC)

	usersRepo := usersLocalStorage.NewUserLocalStorage()
	usersUC := usersUseCase.NewUsersUseCase(usersRepo)

	moviesRepo := moviesLocalStorage.NewMovieLocalStorage()
	moviesUC := moviesUseCase.NewMoviesUseCase(moviesRepo)

	ratingsRepo := ratingsLocalStorage.NewRatingsLocalStorage()
	ratingsUC := ratingsUseCase.NewRatingsUseCase(ratingsRepo)

	authMiddleware := middleware.NewAuthMiddleware(usersUC, sessionsDL)

	return &App{
		usersUC:        usersUC,
		moviesUC:       moviesUC,
		ratingsUC:      ratingsUC,
		sessions:       sessionsDL,
		authMiddleware: authMiddleware,
	}
}

func (app *App) Run(port string) error {
	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "http://89.208.198.186:3000"}
	config.AllowCredentials = true
	router.Use(cors.New(config))

	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))

	router.Static("/avatars", _const.AvatarsFileDir)

	router.Use(gin.Recovery())

	usersHttp.RegisterHttpEndpoints(router, app.usersUC, app.sessions, app.authMiddleware)
	moviesHttp.RegisterHttpEndpoints(router, app.moviesUC)
	ratingsHttp.RegisterHttpEndpoints(router, app.ratingsUC, app.authMiddleware)

	app.server = &http.Server{
		Addr:           ":" + port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		err := app.server.ListenAndServe()
		if err != nil {
			log.Fatal("Failed to listen and serve: ", err)
		}
	}()

	// using graceful shutdown

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return app.server.Shutdown(ctx)
}
