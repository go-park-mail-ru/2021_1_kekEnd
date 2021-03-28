package server

import (
    "context"
    "fmt"
    "github.com/jackc/pgx/v4/pgxpool"
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    "github.com/go-park-mail-ru/2021_1_kekEnd/internal/middleware"
    "github.com/go-park-mail-ru/2021_1_kekEnd/internal/movies"
    moviesHttp "github.com/go-park-mail-ru/2021_1_kekEnd/internal/movies/delivery/http"
    moviesLocalStorage "github.com/go-park-mail-ru/2021_1_kekEnd/internal/movies/repository/localstorage"
    moviesUseCase "github.com/go-park-mail-ru/2021_1_kekEnd/internal/movies/usecase"
    "github.com/go-park-mail-ru/2021_1_kekEnd/internal/reviews"
    reviewsHttp "github.com/go-park-mail-ru/2021_1_kekEnd/internal/reviews/delivery/http"
    reviewsLocalStorage "github.com/go-park-mail-ru/2021_1_kekEnd/internal/reviews/repository/localstorage"
    reviewsUseCase "github.com/go-park-mail-ru/2021_1_kekEnd/internal/reviews/usecase"
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
    reviewsUC      reviews.UseCase
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


    // TO DO Сделать конфиг-файл и считывать его, например, с помощью viper
    connStr := "postgres://mdb:mdb@localhost:5432/mdb"
    dbpool, err := pgxpool.Connect(context.Background(), connStr)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
        os.Exit(1)
    }

    usersRepo := usersLocalStorage.NewUserRepository(dbpool)
    usersUC := usersUseCase.NewUsersUseCase(usersRepo)

    moviesRepo := moviesLocalStorage.NewMovieRepository(dbpool)
    moviesUC := moviesUseCase.NewMoviesUseCase(moviesRepo)

    reviewsRepo := reviewsLocalStorage.NewReviewRepository(dbpool)
    reviewsUC := reviewsUseCase.NewReviewsUseCase(reviewsRepo, usersRepo)

    authMiddleware := middleware.NewAuthMiddleware(usersUC, sessionsDL)

    return &App{
        usersUC:        usersUC,
        moviesUC:       moviesUC,
        sessions:       sessionsDL,
        reviewsUC:      reviewsUC,
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
    reviewsHttp.RegisterHttpEndpoints(router, app.reviewsUC, app.usersUC, app.authMiddleware)

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
