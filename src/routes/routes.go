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
	"time"

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
	msgStr := ""
	hideFiles := ""
	for _, file := range files {

		//Default value is folder
		fileLink := fp.Join("files", path, file.Name)
		download := ""
		fileType := "folder"

		//If it is not a folder, set to file
		if file.IsDir != 0 {
			fileLink = fp.Join("download", file.DownloadName)
			download = "download='" + file.Name + "'"
			fileType = "description"
		}

		//Transform files into html
		htmlStr = htmlStr + "<tr class='item'>" +
			"<td class='id'>" + fmt.Sprint(file.ID) + "</td>" +
			"<td class='name'>" +
			"<a href='/" + fileLink + "' " + download + ">" +
			"<span class='material-icons'>" + fileType + "</span>" +
			"<div>" + file.Name + "</div>" +
			"</a>" +
			"</td>" +
			"</tr>"

	}

	if htmlStr == "" {
		msgStr = "<h1>Não foram encontrados arquivos</h1>"
		hideFiles = "hidden"
	}

	return c.Render("files", fiber.Map{
		"Files":       template.HTML(htmlStr),
		"MsgNotFound": template.HTML(msgStr),
		"FilesHidden": hideFiles,
		"Path":        path,
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
		rand.Seed(time.Now().Unix())
		downloadName := newName[0] + strconv.Itoa(rand.Intn(999)) + "." + newName[1]
		database.InsertFile(addpath, filedata.Filename, downloadName, 1, db)

		err = c.SaveFile(filedata, "./uploads/"+downloadName)
		if err != nil {
			log.Fatal(err)
		}
	}

	return c.Redirect("/files/" + addpath)
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

	return c.Redirect("/files")
}

//DownloadFile download a saved file
func DownloadFile(c *fiber.Ctx) error {
	downloadName := c.Params("*")
	file := database.GetFileByDownloadName(database.GetDB("database.db"), downloadName)
	return c.Download("./uploads/" + file.DownloadName)
}

//SearchFiles finds a file to download or a folder to access
func SearchFiles(c *fiber.Ctx) error {
	idStr := c.Query("id", "0")
	name := c.Query("name", "")
	idInt, _ := strconv.Atoi(idStr)

	if idInt != 0 {
		file := database.GetFileByID(database.GetDB("database.db"), int64(idInt))
		if file.IsDir == 0 {
			return c.Redirect("/files/" + file.Path + file.Name)
		}
		return c.Redirect("/download/" + file.DownloadName)
	} else if name != "" {
		files := database.GetFilesByName(database.GetDB("database.db"), name)
		if len(files) == 0 {
			return c.Render("files", fiber.Map{
				"Files":       "",
				"MsgNotFound": template.HTML("<h1>Nenhum arquivo tem esse nome</h1>"),
				"FilesHidden": "hidden",
				"Path":        "",
			})
		}

		htmlStr := ""
		for _, file := range files {

			//Default value is folder
			fileLink := fp.Join("files", file.Path, file.Name)
			download := ""
			fileType := "folder"

			//If it is not a folder, set to file
			if file.IsDir != 0 {
				fileLink = fp.Join("download", file.DownloadName)
				download = "download='" + file.Name + "'"
				fileType = "description"
			}

			//Transform files into html
			htmlStr = htmlStr + "<tr class='item'>" +
				"<td class='id'>" + fmt.Sprint(file.ID) + "</td>" +
				"<td class='name'>" +
				"<a href='/" + fileLink + "' " + download + ">" +
				"<span class='material-icons'>" + fileType + "</span>" +
				"<div>" + file.Name + "</div>" +
				"</a>" +
				"</td>" +
				"</tr>"

		}

		return c.Render("files", fiber.Map{
			"Files":       template.HTML(htmlStr),
			"MsgNotFound": "",
			"FilesHidden": "",
			"Path":        "",
		})
	}
	return c.Render("files", fiber.Map{
		"Files":       "",
		"MsgNotFound": template.HTML("<h1>Coloque um id ou nome para procurar um arquivo</h1>"),
		"FilesHidden": "hidden",
		"Path":        "",
	})
}
