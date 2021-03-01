package main

import (
	"hd-virtual-plus-plus/routes"
	"html/template"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html"
	//"github.com/gofiber/fiber/v2/middleware/basicauth"
)

var templates *template.Template

func main() {
	engine := html.New("./frontend/template", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	//USE ONLY IF ALL FAILS
	// app.Use(basicauth.New(basicauth.Config{
	// 	Users: map[string]string{
	// 		"nome.aluno@p4ed.com": "SenhaDiscreta9876",
	// 	},
	// 	Authorizer: func(user, passwd string) bool {
	// 		return user == "nome.aluno@p4ed.com" || passwd == "SenhaDiscreta9876"
	// 	},
	// }))

	app.Use(logger.New())
	app.Static("/frontend", "./frontend")
	app.Static("/icons", "./frontend/icons")
	app.Static("/download", "./uploads")

	app.Get("/", routes.Index)
	app.Get("/arquivos/*", routes.Files)
	app.Post("/login", routes.Login)
	app.Get("/add/*", routes.AddFilesForm)
	app.Post("/add/", routes.SaveFiles)

	app.Listen(":3000")
}
