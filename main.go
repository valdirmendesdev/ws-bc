package main

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"github.com/valdirmendesdev/ws-bc/adapters/http/handlers/bc"
	"github.com/valdirmendesdev/ws-bc/config"
	"github.com/valdirmendesdev/ws-bc/docs"
)

// @title         API Serviços do banco central do Brasil
// @version       1.0
// @description   Documentação técnica para utilização dos serviços do banco central do Brasil
//
// @contact.name  Orbit Team
// @contact.email produtos.cloud@seidor.com.br
//
// @schemes       https
//
// @BasePath      /
func main() {

	app := fiber.New()

	// Middleware
	app.Use(recover.New())
	app.Use(cors.New())
	app.Use(logger.New())

	sc := config.New(os.Getenv("HOST"), os.Getenv("PORT"))

	fillAPIDocSettings(sc)

	app.Get("/series/:series_number", bc.Series())
	app.Get("/series/:series_number/latest/:quantity", bc.SeriesUltimos())
	app.Get("/internal-docs/*", swagger.HandlerDefault)
	app.Listen(fmt.Sprintf(":%s", sc.Port()))
}

func fillAPIDocSettings(sc *config.ServiceConfig) {

	docs.SwaggerInfo.Host = sc.FullHost()
}
