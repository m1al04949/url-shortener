package main

import (
	"log"

	"github.com/m1al04949/url-shortener/internal/app"
)

func main() {

	if err := app.RunServer(); err != nil {
		log.Fatal(err)
	}

}
