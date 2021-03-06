package main

import (
	"github.com/go-park-mail-ru/2021_1_kekEnd/server"
	"log"
)

const port string = "8080"

func main() {
	app := server.NewApp()

	if err := app.Run(port); err != nil {
		log.Fatal(err)
	}
}
