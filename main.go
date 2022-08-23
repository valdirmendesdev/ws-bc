package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/valdirmendesdev/ws-bc/adapters/http/handlers/bc"
)

func main() {
	app := fiber.New()

	app.Get("/series/:series_number", bc.Series())
	app.Get("/series/:series_number/latest/:quantity", bc.SeriesUltimos())
	app.Listen(":8080")
}
