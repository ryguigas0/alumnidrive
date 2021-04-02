package userbase

import (
	"database/sql"
	"log"
	"math/rand"
	"time"

	//Imports the sqlite3 drivers
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

//FileModel is model of the columns of the database
type UserModel struct {
	ID     int64
	Name   string
	Email  string
	Psswd  string
	Logged int
	Token  string
}

//GetDB Finds the database in the path and creates the needed tables
func GetDB(path string) *sql.DB {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := db.Prepare("create table if not exists users (id integer primary key unique, name text, email text unique, passwd text, logged int, token text unique)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func InsertNewUser(name, email, psswd string, db *sql.DB) {
	encrPswd, err := bcrypt.GenerateFromPassword([]byte(psswd), 14)
	if err != nil {
		log.Fatalln("Could not encrypt password: " + err.Error())
	}
	stmt, err := db.Prepare("insert into users (id, name, email, passwd) values (?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	rand.Seed(time.Now().UnixNano())
	id := rand.Int63()
	_, err = stmt.Exec(id, name, email, encrPswd)
	if err != nil {
		log.Fatal(err)
	}
}

func FindUserByEmail(email string, db *sql.DB) UserModel {
	row := db.QueryRow("select * from users where email = ?", email)
	if row.Err() != nil {
		log.Fatalf("NO USER WITH EMAIL: %s: %s\n", email, row.Err().Error())
	}
	var usr UserModel
	row.Scan(&usr.ID, &usr.Name, &usr.Email, &usr.Psswd, &usr.Logged, &usr.Token)
	return usr
}
