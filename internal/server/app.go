package server

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/auth"
	authHttp "github.com/go-park-mail-ru/2021_1_kekEnd/internal/auth/delivery/http"
	authLocalStorage "github.com/go-park-mail-ru/2021_1_kekEnd/internal/auth/repository/localstorage"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/auth/usecase"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type App struct {
	server *http.Server
	authUC auth.UseCase
}

func NewApp() *App {
	repo := authLocalStorage.NewUserLocalStorage()

	return &App{
		authUC: usecase.NewAuthUseCase(repo),
	}
}

func (app *App) Run(port string) error {
	router := gin.Default()

	authHttp.RegisterHttpEndpoints(router, app.authUC)

	app.server = &http.Server{
		Addr:           ":" + port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := app.server.ListenAndServe(); err != nil {
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
