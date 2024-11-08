package main

import (
	"github.com/devops-geomais/radarscrapper/handlers"
	"github.com/gofiber/fiber/v2"

)

func setupRoutes(app *fiber.App) {
	app.Get("/", handlers.Home)
	app.Get("/versaoesus", handlers.VersaoEsus)
	app.Get("/verificaversao", handlers.VerificaVersao)
}