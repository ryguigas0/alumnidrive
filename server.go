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
		//MAX FILE SIZE: 100MB
		BodyLimit: 100 * 1024 * 1024,
		Prefork:   true,
	})

	app.Use(logger.New())
	app.Use(recover.New())

	app.Static("static", "src/views/css")
	app.Static("static", "src/views/js")
	app.Static("static", "src/views/icons")

	app.Get("/", routes.Index)
	app.Get("/files/*", routes.Files)
	app.Get("/download/*", routes.DownloadFile)
	app.Get("/add/*", routes.AddFilesForm)
	app.Get("/search", routes.SearchFiles)
	app.Post("/login", routes.Login)
	app.Post("/add/", routes.SaveFiles)
	app.Post("/remove", routes.DeleteFile)

	app.Listen(":3000")
}
