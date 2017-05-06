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

	res.Header().Set("Content-Type", "application/json")

	if req.Method != "POST" {
		fmt.Println("Not a Post Request")
		initJson(0, "Not a Post Request")
		json.NewEncoder(res).Encode(jsonMap)
		return
	}

	var user model.UserLogin

	body, err := ioutil.ReadAll(req.Body)

	if err != nil {
		fmt.Println("Error in Reading")
		initJson(0, "Invalid Data")
		json.NewEncoder(res).Encode(jsonMap)
		return
	}

	err = json.Unmarshal(body, &user)

	if err != nil {
		fmt.Println("Invalid Keys")
		initJson(0, "Invalid Keys")
		json.NewEncoder(res).Encode(jsonMap)
		return
	}

	exists := databaseHelper.CheckUsernameExists(user.Username)

	if !exists {
		fmt.Println("Username doesnot exists")
		initJson(0, "Username doesnot exists")
		json.NewEncoder(res).Encode(jsonMap)
		return
	}

	//check if already logged in
	var token string
	err, token = databaseHelper.GetToken(user.Username)

	if err != nil {
		initJson(0, "Internal Error")
		json.NewEncoder(res).Encode(jsonMap)
		return
	}

	var tokenExist bool
	err, tokenExist = databaseHelper.CheckTokenInDb(token)

	if tokenExist || err != nil {
		fmt.Println("User Already Logged In")
		initJson(0, "User Already Logged In")
		json.NewEncoder(res).Encode(jsonMap)
		return
	}

	err, token = databaseHelper.CheckPassword(user)

	if err != nil {
		fmt.Println("Wrong Password")
		initJson(0, "Wrong Password")
		json.NewEncoder(res).Encode(jsonMap)
		return
	}

	returnToken := make(map[string]interface{})

	returnToken["status"] = 1
	returnToken["message"] = "User logged In"
	returnToken["token"] = token
	json.NewEncoder(res).Encode(returnToken)
}

//Takes username and token
//Resturns firstname lastname of a user
func UserProfile(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	if req.Method != "POST" {
		initJson(0, "Not a POST request")
		json.NewEncoder(res).Encode(jsonMap)
		return
	}

	body, err := ioutil.ReadAll(req.Body)

	var user model.UserProfile
	if err != nil {
		initJson(0, "Error in Reading")
		json.NewEncoder(res).Encode(jsonMap)
		return
	}

	err = json.Unmarshal(body, &user)

	if err != nil {
		initJson(0, "Error in Parsing")
		json.NewEncoder(res).Encode(jsonMap)
		return
	}

	var exist bool

	err, exist = databaseHelper.CheckTokenInDb(user.Token)

	if err != nil || exist == false {
		fmt.Println("User not present")
		initJson(0, "User not present")
		json.NewEncoder(res).Encode(jsonMap)
		return
	}

	var firstname, lastname string

	firstname, lastname, err = databaseHelper.GetUserDetail(user.Username)

	resultMap := make(map[string]interface{})

	resultMap["status"] = 1
	resultMap["firstname"] = firstname
	resultMap["lastname"] = lastname
	json.NewEncoder(res).Encode(resultMap)
}
