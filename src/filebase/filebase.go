package filebase

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	//Imports the sqlite3 drivers
	_ "github.com/mattn/go-sqlite3"
)

//FileModel is model of the columns of the database
type FileModel struct {
	ID           int64
	Path         string
	Name         string
	IsDir        int
	DownloadName string
}

//GetDB Finds the database in the path and creates the needed tables
func GetDB(path string) *sql.DB {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := db.Prepare("create table if not exists files (id integer primary key unique, path text, name text unique, isdir int, downloadname text, owner_id integer foreign key)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Fatal(err)
	}
	return db
}

//InsertFile inserts a path into the database
func InsertFile(path, name, downloadname string, isdir int, db *sql.DB) {
	stmt, err := db.Prepare("insert into files (id, path, name, isdir, downloadname) values (?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	rand.Seed(time.Now().UnixNano())
	id := rand.Int63()
	_, err = stmt.Exec(id, path, name, isdir, downloadname)
	if err != nil {
		log.Fatal(err)
	}
	return
}

//GetAllFiles reads all rows from table files
func GetAllFiles(db *sql.DB) []FileModel {
	rows, err := db.Query("select * from files")
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}
	var fs []FileModel
	for rows.Next() {
		var f FileModel
		if err := rows.Scan(&f.ID, &f.Path, &f.Name, &f.IsDir, &f.DownloadName); err != nil {
			log.Fatal(err)
		}
		fs = append(fs, f)
	}
	return fs
}

//GetFileByID finds a single file based on ID
func GetFileByID(db *sql.DB, id int64) FileModel {
	row := db.QueryRow("select * from files where id = ?", id)
	if row.Err() != nil {
		log.Output(1, fmt.Sprintf("ERR: NO FILE W/ ID: %v\n", row.Err()))
		return FileModel{}
	}
	var file FileModel
	row.Scan(&file.ID, &file.Path, &file.Name, &file.IsDir, &file.DownloadName)
	return file
}

//GetFilesInPath finds all files that match a certain path from the begining
func GetFilesInPath(db *sql.DB, path string) []FileModel {
	rows, err := db.Query("select * from files where path like '" + path + "'")
	if err != nil {
		log.Output(1, fmt.Sprintf("ERR: NO FILE W/ PATH: %v\n", err))
	}
	var files []FileModel
	for rows.Next() {
		var file FileModel
		rows.Scan(&file.ID, &file.Path, &file.Name, &file.IsDir, &file.DownloadName)
		files = append(files, file)
	}
	return files
}

//GetFilesByName finds all files with that contains the name
func GetFilesByName(db *sql.DB, name string) []FileModel {
	rows, err := db.Query("select * from files where name like '%" + name + "%'")
	if err != nil {
		log.Output(1, fmt.Sprint(err))
		return nil
	}
	var files []FileModel
	for rows.Next() {
		var file FileModel
		rows.Scan(&file.ID, &file.Path, &file.Name, &file.IsDir, &file.DownloadName)
		files = append(files, file)
	}
	return files
}

//DeleteFile Deletes a file entry based on id or matching path
func DeleteFile(db *sql.DB, id int64) {
	stmt, err := db.Prepare("delete from files where id = ?")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(id)
	if err != nil {
		log.Fatal(err)
	}
	return
}
