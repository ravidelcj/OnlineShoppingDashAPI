package user

import (
	"databaseHelper"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"model"
	"net/http"
)

var jsonMap = make(map[string]interface{})

//A global json object which is sent in case error occurs
//it takes the status and message of the error
//status 0 indicates error
func initJson(status int, message string) {
	jsonMap["status"] = status
	jsonMap["message"] = message
}

//Check data received is not empty
func checkRegisterJson(user model.User) bool {

	if user.FirstName == "" || user.LastName == "" || user.Username == "" || user.Password == "" {
		return false
	}
	return true
}

func encryptPassword(pass string) (string, error) {
	bytePassword, errPasswordHash := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.MinCost)
	return string(bytePassword), errPasswordHash
}

//Register a new User to database
func RegisterUser(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	if req.Method != "POST" {
		fmt.Println("Not a Post Request")
		initJson(0, "Not a Post Request")
		json.NewEncoder(res).Encode(jsonMap)
		return
	}

	body, err := ioutil.ReadAll(req.Body)
	var user model.User

	err = json.Unmarshal(body, &user)
	if err != nil {
		fmt.Println(err)
		initJson(0, "Incorrect keys")
		json.NewEncoder(res).Encode(jsonMap)
		return
	}

	isTrue := checkRegisterJson(user)

	if isTrue == false {
		fmt.Println("Wrong or Incomplete Data")
		initJson(0, "Wrong or Incomplete Data")
		json.NewEncoder(res).Encode(jsonMap)
		return
	}

	usernameExists := databaseHelper.CheckUsernameExists(user.Username)

	if usernameExists {
		fmt.Println("Username Exists")
		initJson(0, "Username Exists")
		json.NewEncoder(res).Encode(jsonMap)
		return
	}

	encryptPass, errEncypt := encryptPassword(user.Password)
	user.Password = encryptPass

	if errEncypt != nil {
		fmt.Println("Error in encryption")
		initJson(0, "Internal Server Error")
	}

	err = databaseHelper.RegisterUserToDb(user)

	if err != nil {
		fmt.Println(err)
		initJson(0, "Error in registering")
		json.NewEncoder(res).Encode(jsonMap)
		return
	}

	initJson(1, "User Registered Successfully")
	json.NewEncoder(res).Encode(jsonMap)
	fmt.Println("User Registered")
}

//Authenticate and login a user and generate a token which is stored in the token table which helps to
//know that the user is signed in
func LoginUser(res http.ResponseWriter, req *http.Request) {

}

//Resturns profile of a user
func UserProfile(res http.ResponseWriter, req *http.Request) {

}
