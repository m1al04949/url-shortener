package main

import (
	"log"

	"github.com/m1al04949/url-shortener/internal/app"
	_ "github.com/m1al04949/url-shortener/internal/lib/api/response"
)

// @title URL Shortener API
// @version 1.0
// @description API сервиса для сокращения URL
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8081
// @BasePath ./cmd/url-shortener
func main() {

	if err := app.RunServer(); err != nil {
		log.Fatal(err)
	}

}
