package main

import (
	"hd-virtual++/filefinder"
	"html/template"
	"log"
	fp "path/filepath"
	"strings"

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

	app.Get("/", func(c *fiber.Ctx) error {
		if c.FormValue("auth") == "false" {
			log.Output(1, "USER OR PASSWORD WRONG")
		}
		return c.SendFile("./frontend/html/login.html", false)
	})

	app.Get("/arquivos/*", func(c *fiber.Ctx) error {
		filePath := c.Params("*")
		fileNames, err := filefinder.FindFiles("uploads/" + filePath)
		if err != nil {
			log.Fatalf("ERROR: FILE FINDER: %v\n", err)
			return c.SendFile("frontend/html/fileNotFound.html")
		}
		htmlStr := ""
		for _, fileName := range fileNames {

			//Default file is a folder
			fileLink := fp.Join("arquivos", filePath, fileName)
			download := ""
			fileType := "folder"

			//If it has an extension it is a file
			if strings.Contains(fileName, ".") {
				fileLink = fp.Join("download", filePath, fileName)
				download = "download='" + fileName + "'"
				fileType = "description"
			}

			//Transform files into html
			htmlStr = htmlStr + "<a href='/" + fileLink + "' " + download + " class='item'>" +
				"<span class='material-icons'>" +
				fileType +
				"</span>" +
				"<div class='name'>" +
				fileName +
				"</div>" +
				"</a>"

		}
		html := template.HTML(htmlStr)
		return c.Render("files", fiber.Map{
			"Files": html,
		})
	})

	app.Post("/login", func(c *fiber.Ctx) error {
		username := c.FormValue("username")
		passwd := c.FormValue("password")

		log.Output(0, username)
		log.Output(0, passwd)
		// if username != "john" || passwd != "doe" {
		// 	return c.Redirect("../", fiber.StatusUnauthorized)
		// }

		return c.Redirect("/arquivos")
	})

	app.Listen(":3000")
}
