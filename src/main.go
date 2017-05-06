package main

import (
	"databaseHelper"
	"net/http"
	"user"
	//"fmt"
	//"model"
)

//User Related Task

func main() {

	err := databaseHelper.InitDatabase()
	if err != nil {
		return
	}
	defer databaseHelper.Db.Close()

	//User Related URL

	// /user/register
	http.HandleFunc("/user/register", user.RegisterUser)
	// /user/login
	http.HandleFunc("/user/login", user.LoginUser)
	// /user/profile
	http.HandleFunc("/user/profile", user.UserProfile)

	http.ListenAndServe(":8008", nil)

}
