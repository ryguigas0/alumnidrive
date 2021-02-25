package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Static("/", "../frontend")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendFile("../frontend/html/login.html", false)
	})

	app.Get("/arquivos", func(c *fiber.Ctx) error {
		return c.SendFile("../frontend/html/arquivos.html")
	})

	app.Post("/login", func(c *fiber.Ctx) error {

		log.Output(0, c.FormValue("username"))
		//log.Output(0, c.FormValue("password"))

		return c.Redirect("/arquivos")
	})

	app.Listen(":3000")
}
