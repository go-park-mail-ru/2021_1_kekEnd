package main

import (
	"log"

	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/server"
	constants "github.com/go-park-mail-ru/2021_1_kekEnd/pkg/const"
)

func main() {
	app := server.NewApp()

	if err := app.Run(constants.Port); err != nil {
		log.Fatal(err)
	}
}
