package databaseHelper

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"model"
	//"golang.org/x/crypto/bcrypt"
)

//Global Database variable
var Db *sql.DB

//Initialise Database
func InitDatabase() error {

	var err error
	Db, err = sql.Open("mysql", "root:admin@/online_store")

	if err != nil {
		fmt.Println("Error Connecting to Database : ", err)
		return err
	}

	err = Db.Ping()

	if err != nil {
		fmt.Println("Error Contacting Database : ", err)
		return err
	}
	return nil
}

func CheckUsernameExists(username string) bool {
	query := "select exists(select 1 from users where username = '" + username + "');"

	var exists bool

	err := Db.QueryRow(query).Scan(&exists)

	if err != nil {
		fmt.Println(err)
		return true
	}
	return exists
}

func RegisterUserToDb(user model.User) error {
	stmt, err := Db.Prepare("Insert into users set username = ? , password = ? , firstname = ? , lastname = ? ")

	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = stmt.Exec(user.Username, user.Password, user.FirstName, user.LastName)

	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
