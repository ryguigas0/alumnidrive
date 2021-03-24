package routes

import (
	"fmt"
	"hd-virtual-plus-plus/src/database"
	"html/template"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	fp "path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

var errMap map[string]string = map[string]string{
	"dir":  "The directory has no name",
	"file": "The file is empty",
	"name": "This name is already taken",
}

//Index Load login page, with or without warning of wrong user or password
func Index(c *fiber.Ctx) error {
	//c.Cookie(&fiber.Cookie{Name: "darktheme", Value: "1"})
	return c.Render("login", fiber.Map{})
}

//AddFilesForm Add files or folders to a form and send
func AddFilesForm(c *fiber.Ctx) error {
	pathDir := c.Params("*")
	pathName := pathDir
	if len(pathDir) == 0 {
		pathDir = ""
		pathName = "pasta inicial"
	}
	errMsg := c.Query("err", "")
	return c.Render("addFileForm", fiber.Map{
		"PathName": pathName,
		"PathDir":  pathDir,
		"ErrorMsg": errMap[errMsg],
	})
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

//filesToHTML turns a file list into html
func filesToHTML(files []database.FileModel) (htmlStr, msgStr, hideFiles string) {
	for _, file := range files {

		//Default value is folder
		fileLink := fp.Join("files", file.Path, file.Name)
		download := ""
		fileType := "folder"

		//If it is not a folder, set to file
		if file.IsDir != 0 {
			fileLink = fp.Join("download", fmt.Sprint(file.ID))
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
			"<td>" +
			"<form method='post' action='/remove'>" +
			"<input type='hidden' name='id' value='" + fmt.Sprint(file.ID) + "'>" +
			"<span class='material-icons' onclick='this.parentNode.submit();'> do_not_disturb_on </span>" +
			"</form>" +
			"</td>" +
			"</tr>"
	}

	if htmlStr == "" {
		msgStr = "<h1>Não foram encontrados arquivos</h1>"
		hideFiles = "hidden"
	}

	return
}

//Files return all files within a path
func Files(c *fiber.Ctx) error {
	path := c.Params("*")
	outPathHTMLStr := ""
	if path != "" {
		outPath := strings.Split(path, "/")
		outPathStr := filepath.Join(outPath[:len(outPath)-1]...)
		outPathHTMLStr = "<a class='out-path-redir' href='/files/" + outPathStr + "'>" +
			"<span class='material-icons'>" +
			"arrow_back" +
			"</span>" +
			"</a>"
	}

	files := database.GetFilesInPath(database.GetDB(("database.db")), path)
	htmlStr, msgStr, hideFiles := filesToHTML(files)
	return c.Render("files", fiber.Map{
		"Files":       template.HTML(htmlStr),
		"Mensage":     template.HTML(msgStr),
		"FilesHidden": hideFiles,
		"Path":        path,
		"OutPath":     template.HTML(outPathHTMLStr),
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
			return c.Redirect("/add/" + addpath + "?err=dir")
		}
		if filesWithName := database.GetFilesByName(db, dirname); len(filesWithName) != 0 {
			return c.Redirect("/add/" + addpath + "?err=name")
		}
		database.InsertFile(addpath, dirname, dirname, 0, db)
	} else {
		filedata, err := c.FormFile("filedata")
		if err != nil {
			return c.Redirect("/add/" + addpath + "?err=file")
		}
		newName := strings.Split(strings.ReplaceAll(filedata.Filename, " ", "_"), ".")
		if filesWithName := database.GetFilesByName(db, strings.Join([]string{newName[0], newName[1]}, ".")); len(filesWithName) != 0 {
			return c.Redirect("/add/" + addpath + "?err=name")
		}
		rand.Seed(time.Now().Unix())
		downloadName := newName[0] + strconv.Itoa(rand.Intn(999)) + "." + newName[1]
		database.InsertFile(addpath, strings.Join([]string{newName[0], newName[1]}, "."), downloadName, 1, db)

		err = c.SaveFile(filedata, "./uploads/"+downloadName)
		if err != nil {
			log.Fatal(err)
		}
	}

	return c.Redirect("/files/" + addpath)
}

func deleteSavedFile(file database.FileModel) {
	db := database.GetDB("database.db")
	if file.IsDir != 0 {
		err := os.Remove(filepath.Join(".", "uploads", file.DownloadName))
		if err != nil {
			log.Output(1, fmt.Sprint("Can't delete file: ", err))
		}
	} else {
		files := database.GetFilesInPath(db, filepath.Join(file.Path, file.Name))
		for _, file := range files {
			deleteSavedFile(file)
		}
	}
	database.DeleteFile(db, file.ID)
}

//DeleteFile deletes a file with its id
func DeleteFile(c *fiber.Ctx) error {
	idStr := c.FormValue("id", "0")
	idInt, _ := strconv.Atoi(idStr)
	if idInt != 0 {
		db := database.GetDB("database.db")
		file := database.GetFileByID(db, int64(idInt))

		deleteSavedFile(file)

		msgStr := "<h1>O arquivo " + idStr + " foi deletado!</h1>"
		if file.IsDir == 0 {
			msgStr = "<h1>A pasta " + idStr + " e seus conteúdos foram deletados!</h1>"
		}

		return c.Render("fileDeleted", fiber.Map{
			"Mensage": template.HTML(msgStr),
			"Path":    file.Path,
		})
	}
	return c.Render("files", fiber.Map{
		"Mensage": template.HTML("<h1>ID inválido!</h1>"),
	})
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
		return c.Redirect("/download/" + fmt.Sprint(file.ID))
	} else if name != "" {
		files := database.GetFilesByName(database.GetDB("database.db"), name)
		htmlStr, msgStr, hideFiles := filesToHTML(files)
		return c.Render("files", fiber.Map{
			"Files":       template.HTML(htmlStr),
			"Mensage":     template.HTML(msgStr),
			"FilesHidden": hideFiles,
			"Path":        "",
		})
	}
	return c.Render("files", fiber.Map{
		"Files":       "",
		"Mensage":     template.HTML("<h1>Coloque um id ou nome para procurar um arquivo</h1>"),
		"FilesHidden": "hidden",
		"Path":        "",
	})
}

//DownloadFile download a saved file
func DownloadFile(c *fiber.Ctx) error {
	idStr := c.Params("*", "0")
	idInt, _ := strconv.Atoi(idStr)
	file := database.GetFileByID(database.GetDB("database.db"), int64(idInt))
	return c.Download("./uploads/" + file.DownloadName)
}
