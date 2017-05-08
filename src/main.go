package main

import (
	"databaseHelper"
	"net/http"
	"product"
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

	//Add product to database
	http.HandleFunc("/product/add", product.AddProd)
	//Delete a particular product
	// http.HandleFunc("/product/delete", product.DeleteProd)
	// //Edit a particular product
	// http.HandleFunc("/product/edit", product.EditProd)
	// //Search Products
	// http.HandleFunc("product/search", product.SearchProd)

	http.ListenAndServe(":8008", nil)

}
