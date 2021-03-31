package main

import (
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/server"
	_const "github.com/go-park-mail-ru/2021_1_kekEnd/pkg/const"
	"log"
)

func main() {
	app := server.NewApp()

	if err := app.Run(_const.Port); err != nil {
		log.Fatal(err)
	}
}
