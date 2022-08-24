package main

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/valdirmendesdev/ws-bc/adapters/http/handlers/bc"
)

func main() {
	app := fiber.New()
	app.Use(logger.New())
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	} else {
		port = fmt.Sprintf(":%s", port)
	}

	app.Get("/series/:series_number", bc.Series())
	app.Get("/series/:series_number/latest/:quantity", bc.SeriesUltimos())
	app.Listen(port)
}
