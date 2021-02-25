package main

import (
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

	app.Listen(":3000")
}
