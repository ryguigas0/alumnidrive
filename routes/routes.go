package routes

import (
	"fmt"
	"hd-virtual-plus-plus/filefinder"
	"html/template"
	"log"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
)

//Index Load login page, with or without warning of wrong user or password
func Index(c *fiber.Ctx) error {
	if c.FormValue("auth") == "false" {
		log.Output(1, "USER OR PASSWORD WRONG")
	}
	return c.SendFile("./frontend/html/login.html", false)
}

//Files serve files to user based on url path
func Files(c *fiber.Ctx) error {
	filePath := c.Params("*")
	fileNames, err := filefinder.FindFiles("uploads/" + filePath)
	if err != nil {
		log.Fatalf("ERROR: FILE FINDER: %v\n", err)
		return c.SendFile("frontend/html/fileNotFound.html")
	}
	htmlStr := ""
	for _, fileName := range fileNames {

		//Default file is a folder
		fileLink := filepath.Join("arquivos", filePath, fileName)
		download := ""
		fileType := "folder"

		//If it has an extension it is a file
		if strings.Contains(fileName, ".") {
			fileLink = filepath.Join("download", filePath, fileName)
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
		"Path":  filePath,
	})
}

//AddFilesForm Add files or folders to a form and send
func AddFilesForm(c *fiber.Ctx) error {
	local := c.Params("*")
	if len(local) == 0 {
		local = "pasta inicial"
	}
	return c.Render("addFileForm", fiber.Map{
		"Path": local,
	})
}

//SaveFiles Save files from form and tag them with their id
func SaveFiles(c *fiber.Ctx) error {
	filename := c.FormValue("filename")
	filetype := c.FormValue("filetype")
	filepath := c.FormValue("filepath")
	filedata := c.FormValue("filedata")
	return c.Send([]byte(fmt.Sprintf("%v %v %v %v", filename, filetype, filepath, filedata)))
}

//Login login the user and give access to uploaded files
func Login(c *fiber.Ctx) error {
	username := c.FormValue("username")
	passwd := c.FormValue("password")

	log.Output(0, username)
	log.Output(0, passwd)
	// if username != "john" || passwd != "doe" {
	// 	return c.Redirect("../", fiber.StatusUnauthorized)
	// }

	return c.Redirect("/arquivos")
}
