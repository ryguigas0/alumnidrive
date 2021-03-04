package main

import (
	"hd-virtual-plus-plus/src/routes"
	"html/template"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html"
	//"github.com/gofiber/fiber/v2/middleware/basicauth"
)

var templates *template.Template

func main() {
	engine := html.New("./src/views/pages", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Use(logger.New())
	app.Use(recover.New())

	app.Static("frontend", "src/views/css")
	app.Static("frontend", "src/views/js")
	app.Static("frontend", "src/views/icons")

	app.Get("/", routes.Index)
	app.Get("/files/*", routes.Files)
	app.Get("/download/*", routes.DownloadFile)
	app.Post("/login", routes.Login)
	app.Get("/add/*", routes.AddFilesForm)
	app.Post("/add/", routes.SaveFiles)
	app.Get("/search", routes.SearchFiles)

	app.Listen(":3000")
}
