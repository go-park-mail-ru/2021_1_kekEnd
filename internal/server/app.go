package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/actors"
	actorsHttp "github.com/go-park-mail-ru/2021_1_kekEnd/internal/actors/delivery/http"
	actorsDBStorage "github.com/go-park-mail-ru/2021_1_kekEnd/internal/actors/repository/dbstorage"
	actorsUseCase "github.com/go-park-mail-ru/2021_1_kekEnd/internal/actors/usecase"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/logger"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/middleware"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/movies"
	moviesHttp "github.com/go-park-mail-ru/2021_1_kekEnd/internal/movies/delivery/http"
	moviesDBStorage "github.com/go-park-mail-ru/2021_1_kekEnd/internal/movies/repository/dbstorage"
	moviesUseCase "github.com/go-park-mail-ru/2021_1_kekEnd/internal/movies/usecase"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/playlists"
	playlistsHttp "github.com/go-park-mail-ru/2021_1_kekEnd/internal/playlists/delivery"
	playlistsRepository "github.com/go-park-mail-ru/2021_1_kekEnd/internal/playlists/repository"
	playlistsUseCase "github.com/go-park-mail-ru/2021_1_kekEnd/internal/playlists/usecase"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/proto"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/ratings"
	ratingsHttp "github.com/go-park-mail-ru/2021_1_kekEnd/internal/ratings/delivery"
	ratingsDBStorage "github.com/go-park-mail-ru/2021_1_kekEnd/internal/ratings/repository/dbstorage"
	ratingsUseCase "github.com/go-park-mail-ru/2021_1_kekEnd/internal/ratings/usecase"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/reviews"
	reviewsHttp "github.com/go-park-mail-ru/2021_1_kekEnd/internal/reviews/delivery/http"
	reviewsDBStorage "github.com/go-park-mail-ru/2021_1_kekEnd/internal/reviews/repository/dbstorage"
	reviewsUseCase "github.com/go-park-mail-ru/2021_1_kekEnd/internal/reviews/usecase"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/search"
	searchHttp "github.com/go-park-mail-ru/2021_1_kekEnd/internal/search/delivery/http"
	searchUseCase "github.com/go-park-mail-ru/2021_1_kekEnd/internal/search/usecase"
	sessionsDelivery "github.com/go-park-mail-ru/2021_1_kekEnd/internal/sessions/delivery"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/users"
	usersHttp "github.com/go-park-mail-ru/2021_1_kekEnd/internal/users/delivery/http"
	usersDBStorage "github.com/go-park-mail-ru/2021_1_kekEnd/internal/users/repository/dbstorage"
	usersUseCase "github.com/go-park-mail-ru/2021_1_kekEnd/internal/users/usecase"
	constants "github.com/go-park-mail-ru/2021_1_kekEnd/pkg/const"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

// App структура главного приложения
type App struct {
	server         *http.Server
	usersUC        users.UseCase
	actorsUC       actors.UseCase
	moviesUC       movies.UseCase
	ratingsUC      ratings.UseCase
	reviewsUC      reviews.UseCase
	playlistsUC    playlists.UseCase
	searchUC       search.UseCase
	authMiddleware middleware.Auth
	csrfMiddleware middleware.Csrf
	logger         *logger.Logger
	sessionsDL     *sessionsDelivery.AuthClient
	sessionsConn   *grpc.ClientConn
	fileServer     proto.FileServerHandlerClient
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

// NewApp инициализация приложения
func NewApp() *App {
	accessLogger := logger.NewAccessLogger()

	connStr, connected := os.LookupEnv("DB_CONNECT")
	if !connected {
		fmt.Println(os.Getwd())
		log.Fatal("Failed to read DB connection data")
	}
	dbpool, err := pgxpool.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	sessionsGrpcConn, err := grpc.Dial(fmt.Sprintf("localhost:%s", constants.AuthPort), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Unable to connect to grpc auth server: %v\n", err)
	}
	sessionsDL := sessionsDelivery.NewAuthClient(sessionsGrpcConn)

	fileServerGrpcConn, err := grpc.Dial(fmt.Sprintf("localhost:%s", constants.FileServerPort), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Unable to connect to grpc file server: %v\n", err)
	}
	fileServerService := proto.NewFileServerHandlerClient(fileServerGrpcConn)

	usersRepo := usersDBStorage.NewUserRepository(dbpool)
	reviewsRepo := reviewsDBStorage.NewReviewRepository(dbpool)
	moviesRepo := moviesDBStorage.NewMovieRepository(dbpool)
	actorsRepo := actorsDBStorage.NewActorRepository(dbpool)
	ratingsRepo := ratingsDBStorage.NewRatingsRepository(dbpool)

	usersUC := usersUseCase.NewUsersUseCase(usersRepo, reviewsRepo, ratingsRepo, actorsRepo)
	moviesUC := moviesUseCase.NewMoviesUseCase(moviesRepo, usersRepo)
	actorsUC := actorsUseCase.NewActorsUseCase(actorsRepo)
	reviewsUC := reviewsUseCase.NewReviewsUseCase(reviewsRepo, usersRepo)
	ratingsUC := ratingsUseCase.NewRatingsUseCase(ratingsRepo)
	searchUC := searchUseCase.NewSearchUseCase(usersRepo, moviesRepo, actorsRepo)

	playlistsRepo := playlistsRepository.NewPlaylistsRepository(dbpool)
	playlistsUC := playlistsUseCase.NewPlaylistsUseCase(playlistsRepo)

	authMiddleware := middleware.NewAuthMiddleware(usersUC, sessionsDL)
	csrfMiddleware := middleware.NewCsrfMiddleware(accessLogger)

	return &App{
		usersUC:        usersUC,
		actorsUC:       actorsUC,
		moviesUC:       moviesUC,
		ratingsUC:      ratingsUC,
		reviewsUC:      reviewsUC,
		playlistsUC:    playlistsUC,
		searchUC:       searchUC,
		authMiddleware: authMiddleware,
		csrfMiddleware: csrfMiddleware,
		logger:         accessLogger,
		sessionsDL:     sessionsDL,
		sessionsConn:   sessionsGrpcConn,
		fileServer:     fileServerService,
	}
}

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// Run запуск приложения
func (app *App) Run(port string) error {
	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:4000", "https://cinemedia.ru"}
	config.AllowCredentials = true
	router.Use(cors.New(config))
	router.Use(middleware.AccessLogMiddleware(app.logger))

	router.Static("/avatars", constants.AvatarsFileDir)
	router.Static("/posters", constants.PostersFileDir)
	router.Static("/banners", constants.BannersFileDir)
	router.Static("/actors", constants.ActorsFileDir)

	router.Use(gin.Recovery())
	router.GET("/metrics", prometheusHandler())

	api := router.Group("/api")
	v1 := api.Group("/v1")

	usersHttp.RegisterHTTPEndpoints(v1, app.usersUC, app.sessionsDL, app.authMiddleware, app.fileServer, app.logger)
	moviesHttp.RegisterHTTPEndpoints(v1, app.moviesUC, app.authMiddleware, app.logger)
	ratingsHttp.RegisterHTTPEndpoints(v1, app.ratingsUC, app.authMiddleware, app.logger)
	reviewsHttp.RegisterHTTPEndpoints(v1, app.reviewsUC, app.usersUC, app.authMiddleware, app.logger)
	actorsHttp.RegisterHTTPEndpoints(v1, app.actorsUC, app.authMiddleware, app.logger)
	playlistsHttp.RegisterHTTPEndpoints(v1, app.playlistsUC, app.usersUC, app.authMiddleware, app.logger)
	searchHttp.RegisterHTTPEndpoints(v1, app.searchUC, app.logger)

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

	_ = app.sessionsConn.Close()
	return app.server.Shutdown(ctx)
}
