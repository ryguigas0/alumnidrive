package routes

import (
	"fmt"
	"hd-virtual-plus-plus/src/database"
	"hd-virtual-plus-plus/src/fileman"
	"html/template"
	"log"
	"os"
	fp "path/filepath"
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
	filePath := c.Params("*")
	fileModels, err := fileman.FindFiles("uploads/" + filePath)

	if err != nil {
		log.Fatalf("ERROR: FILE FINDER: %v\n", err)
		return c.Render("fileNotFound", fiber.Map{})
	}
	htmlStr := ""
	for _, fileModel := range fileModels {

		//Default value is folder
		filename := fileModel.Name
		fileLink := fp.Join("arquivos", filePath, filename)
		download := ""
		fileType := "folder"
		fileID := database.GetFiles(database.GetDB("database.db"), filename, 0)[0].ID

		//If it is not a folder, set to file
		if !fileModel.IsDir {
			fileLink = fp.Join("download", filePath, filename)
			download = "download='" + filename + "'"
			fileType = "description"
		}

		//Transform files into html
		htmlStr = htmlStr + "<a href='/" + fileLink + "' " + download + " class='item'>" +
			"<span class='material-icons'>" +
			fileType +
			"</span>" +
			"<div class='name'>" +
			filename +
			"</div>" +
			"<div class='id'>" +
			fmt.Sprint(fileID) +
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

		savepath := fp.Join("uploads", addpath, strings.ReplaceAll(dirname, " ", "_"))
		if err := os.Mkdir(savepath, 0755); err != nil {
			log.Fatalf("ERROR: SAVE DIR: %v\n", err)
			return c.Redirect("/add/" + addpath)
		}
		log.Output(1, fmt.Sprintf("%v %v %v", dirname, addtype, savepath))
		database.InsertFile(savepath, db)
	} else {
		filedata, err := c.FormFile("filedata")
		if err != nil {
			log.Fatalf("ERROR: FILE UPLOAD: %v\n", err)
			c.Request().Header.Add("error", "missing-file")
			return c.Redirect("/add/" + addpath)
		}

		savepath := fp.Join("uploads", addpath, strings.ReplaceAll(filedata.Filename, " ", "_"))
		if err = c.SaveFile(filedata, savepath); err != nil {
			log.Fatalf("ERROR: SAVE UPLOAD: %v\n", err)
			return c.Redirect("/add/" + addpath)
		}
		log.Output(1, fmt.Sprintf("%v %v %v", filedata.Filename, addtype, savepath))
		database.InsertFile(savepath, db)
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
