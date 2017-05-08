package databaseHelper

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
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
