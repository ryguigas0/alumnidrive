package routes

import (
	"fmt"
	"hd-virtual-plus-plus/src/database"
	"html/template"
	"log"
	"math/rand"
	fp "path/filepath"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

//Index Load login page, with or without warning of wrong user or password
func Index(c *fiber.Ctx) error {
	if c.FormValue("auth") == "false" {
		log.Output(1, "USER OR PASSWORD WRONG")
	}
	return c.Render("login", fiber.Map{})
}

//Files serve files to user based on url path
func Files(c *fiber.Ctx) error {
	path := c.Params("*")

	db := database.GetDB("database.db")
	defer db.Close()
	files := database.GetFilesInPath(db, path)

	htmlStr := ""
	for _, file := range files {

		//Default value is folder
		fileLink := fp.Join("arquivos", path, file.Name)
		download := ""
		fileType := "folder"

		//If it is not a folder, set to file
		if file.IsDir != 0 {
			fileLink = fp.Join("download", file.DownloadName)
			download = "download='" + file.Name + "'"
			fileType = "description"
		}

		//Transform files into html
		htmlStr = htmlStr + "<a href='/" + fileLink + "' " + download + " class='item'>" +
			"<span class='material-icons'>" +
			fileType +
			"</span>" +
			"<div class='name'>" +
			file.Name +
			"</div>" +
			"<div class='id'>" +
			fmt.Sprint(file.ID) +
			"</div>" +
			"</a>"

	}
	html := template.HTML(htmlStr)
	return c.Render("files", fiber.Map{
		"Files": html,
		"Path":  path,
	})
}

//AddFilesForm Add files or folders to a form and send
func AddFilesForm(c *fiber.Ctx) error {
	pathDir := c.Params("*")
	pathName := pathDir
	if len(pathDir) == 0 {
		pathDir = ""
		pathName = "pasta inicial"
	}
	return c.Render("addFileForm", fiber.Map{
		"PathName": pathName,
		"PathDir":  pathDir,
	})
}

//SaveFiles Save files from form and tag them with their id
func SaveFiles(c *fiber.Ctx) error {
	addpath := c.FormValue("addpath")
	addtype := c.FormValue("addtype")
	db := database.GetDB("database.db")
	defer db.Close()

	if addtype == "dir" {
		dirname := c.FormValue("dirname")

		if dirname == "" {
			c.Request().Header.Add("error", "missing-dirname")
			return c.Redirect("/add/" + addpath)
		}

		database.InsertFile(addpath, dirname, dirname, 0, db)
	} else {
		filedata, err := c.FormFile("filedata")
		if err != nil {
			log.Fatalf("ERROR: FILE UPLOAD: %v\n", err)
			c.Request().Header.Add("error", "missing-file")
			return c.Redirect("/add/" + addpath)
		}
		newName := strings.Split(strings.ReplaceAll(filedata.Filename, " ", "_"), ".")
		downloadName := newName[0] + strconv.Itoa(rand.Intn(999)) + "." + newName[1]
		database.InsertFile(addpath, filedata.Filename, downloadName, 1, db)

		err = c.SaveFile(filedata, "./uploads/"+downloadName)
		if err != nil {
			log.Fatal(err)
		}
	}

	return c.Redirect("/arquivos/" + addpath)
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

//DownloadFile download a saved file
func DownloadFile(c *fiber.Ctx) error {
	downloadName := c.Params("*")
	file := database.GetFileByDownloadName(database.GetDB("database.db"), downloadName)
	return c.Download("./uploads/" + file.DownloadName)
	//return c.Redirect("/arquivos/" + file.Path)
}
