package database

import (
	"database/sql"
	"log"

	//Imports the sqlite3 drivers
	_ "github.com/mattn/go-sqlite3"
)

type File struct {
	ID   int
	Path string
}

//GetDB Finds the database in the path and creates the needed tables
func GetDB(path string) *sql.DB {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := db.Prepare("create table if not exists files (id integer primary key autoincrement, path text unique)")
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
func InsertFile(path string, db *sql.DB) {
	stmt, err := db.Prepare("insert into files (path) values (?)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(path)
	if err != nil {
		log.Fatal(err)
	}
	return
}

//GetFiles reads all files from table
func GetFiles(db *sql.DB, pattern string, id int) []File {
	rows, err := db.Query("select * from files where path LIKE '%"+pattern+"%' or id = ?", id)
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}
	var fs []File
	for rows.Next() {
		var f File
		if err := rows.Scan(&f.ID, &f.Path); err != nil {
			log.Fatal(err)
		}
		fs = append(fs, f)
	}
	return fs
}

//DeleteFile Deletes a file entry based on id or matching path
func DeleteFile(db *sql.DB, pattern string, id int) {
	stmt, err := db.Prepare("delete from files where path LIKE '%" + pattern + "%' or id = ?")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(id)
	if err != nil {
		log.Fatal(err)
	}
	return
}
