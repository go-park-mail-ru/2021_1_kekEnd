package server

import (
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/actors"
	actorsHttp "github.com/go-park-mail-ru/2021_1_kekEnd/internal/actors/delivery/http"
	actorsDBStorage "github.com/go-park-mail-ru/2021_1_kekEnd/internal/actors/repository/dbstorage"
	actorsUseCase "github.com/go-park-mail-ru/2021_1_kekEnd/internal/actors/usecase"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/logger"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/middleware"
	//"github.com/go-park-mail-ru/2021_1_kekEnd/internal/logger"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/movies"
	moviesHttp "github.com/go-park-mail-ru/2021_1_kekEnd/internal/movies/delivery/http"
	moviesDBStorage "github.com/go-park-mail-ru/2021_1_kekEnd/internal/movies/repository/dbstorage"
	moviesUseCase "github.com/go-park-mail-ru/2021_1_kekEnd/internal/movies/usecase"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/ratings"
	ratingsHttp "github.com/go-park-mail-ru/2021_1_kekEnd/internal/ratings/delivery"
	ratingsDBStorage "github.com/go-park-mail-ru/2021_1_kekEnd/internal/ratings/repository/dbstorage"
	ratingsUseCase "github.com/go-park-mail-ru/2021_1_kekEnd/internal/ratings/usecase"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/reviews"
	reviewsHttp "github.com/go-park-mail-ru/2021_1_kekEnd/internal/reviews/delivery/http"
	reviewsDBStorage "github.com/go-park-mail-ru/2021_1_kekEnd/internal/reviews/repository/dbstorage"
	reviewsUseCase "github.com/go-park-mail-ru/2021_1_kekEnd/internal/reviews/usecase"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/sessions"
	sessionsDelivery "github.com/go-park-mail-ru/2021_1_kekEnd/internal/sessions/delivery"
	sessionsRepository "github.com/go-park-mail-ru/2021_1_kekEnd/internal/sessions/repository"
	sessionsUseCase "github.com/go-park-mail-ru/2021_1_kekEnd/internal/sessions/usecase"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/users"
	usersHttp "github.com/go-park-mail-ru/2021_1_kekEnd/internal/users/delivery/http"
	usersDBStorage "github.com/go-park-mail-ru/2021_1_kekEnd/internal/users/repository/dbstorage"
	usersUseCase "github.com/go-park-mail-ru/2021_1_kekEnd/internal/users/usecase"
	_const "github.com/go-park-mail-ru/2021_1_kekEnd/pkg/const"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type App struct {
	server         *http.Server
	usersUC        users.UseCase
	actorsUC 	   actors.UseCase
	moviesUC       movies.UseCase
	ratingsUC      ratings.UseCase
	reviewsUC      reviews.UseCase
	sessions       sessions.Delivery
	authMiddleware middleware.Auth
	logger         *logger.Logger
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
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

	accessLogger := logger.NewAccessLogger()


	sessionsRepo := sessionsRepository.NewRedisRepository(rdb)
	sessionsUC := sessionsUseCase.NewUseCase(sessionsRepo)
	sessionsDL := sessionsDelivery.NewDelivery(sessionsUC, accessLogger)

	connStr, connected := os.LookupEnv("DB_CONNECT")
	if !connected {
		log.Fatal("Failed to read DB connection data", err)
	}
	dbpool, err := pgxpool.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	usersRepo := usersDBStorage.NewUserRepository(dbpool)
	usersUC := usersUseCase.NewUsersUseCase(usersRepo)

	actorsRepo := actorsDBStorage.NewActorRepository(dbpool)
	actorsUC := actorsUseCase.NewActorsUseCase(actorsRepo)

	moviesRepo := moviesDBStorage.NewMovieRepository(dbpool)
	moviesUC := moviesUseCase.NewMoviesUseCase(moviesRepo)

	reviewsRepo := reviewsDBStorage.NewReviewRepository(dbpool)
	reviewsUC := reviewsUseCase.NewReviewsUseCase(reviewsRepo, usersRepo)

	ratingsRepo := ratingsDBStorage.NewRatingsRepository(dbpool)
	ratingsUC := ratingsUseCase.NewRatingsUseCase(ratingsRepo)

	authMiddleware := middleware.NewAuthMiddleware(usersUC, sessionsDL)

	return &App{
		usersUC:        usersUC,
		actorsUC:       actorsUC,
		moviesUC:       moviesUC,
		ratingsUC:      ratingsUC,
		sessions:       sessionsDL,
		reviewsUC:      reviewsUC,
		authMiddleware: authMiddleware,
		logger:         accessLogger,
	}
}

func (app *App) Run(port string) error {
	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "http://89.208.198.186:3000"}
	config.AllowCredentials = true
	router.Use(cors.New(config))
	router.Use(middleware.AccessLogMiddleware(app.logger))

	router.Static("/avatars", _const.AvatarsFileDir)

	router.Use(gin.Recovery())

	usersHttp.RegisterHttpEndpoints(router, app.usersUC, app.sessions, app.authMiddleware, app.logger)
	moviesHttp.RegisterHttpEndpoints(router, app.moviesUC, app.logger)
	ratingsHttp.RegisterHttpEndpoints(router, app.ratingsUC, app.authMiddleware, app.logger)
	reviewsHttp.RegisterHttpEndpoints(router, app.reviewsUC, app.usersUC, app.authMiddleware, app.logger)
	actorsHttp.RegisterHttpEndpoints(router, app.actorsUC, app.authMiddleware, app.logger)

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
