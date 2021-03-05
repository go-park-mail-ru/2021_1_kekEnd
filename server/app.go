package server

import (
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2021_1_kekEnd/auth"
	"github.com/go-park-mail-ru/2021_1_kekEnd/auth/repository/localstorage"
	"github.com/go-park-mail-ru/2021_1_kekEnd/auth/usecase"
	"net/http"
)

type App struct {
	server *http.Server
	authUC auth.UseCase
}

func NewApp() *App {
	repo := localstorage.NewUserLocalStorage()

	return &App{
		authUC: usecase.NewAuthUseCase(repo),
	}
}

func (app *App) Run(port string) error {
	//router := gin.Default()
}
