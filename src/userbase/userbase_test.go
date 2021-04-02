package userbase_test

import (
	"anondrive/src/userbase"
	"fmt"
	"testing"
)

func TestNewAccount(t *testing.T) {
	db := userbase.GetDB("../../database.db")
	usrname := "John Doe"
	usremail := "joe@doe.usr"
	psswd := "joedoe123"
	userbase.InsertNewUser(usrname, usremail, psswd, db)
	if found := userbase.FindUserByEmail(usremail, db); found.ID == 0 {
		t.FailNow()
	} else {
		fmt.Println(found)
	}
}
