package databaseHelper

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"model"
	"time"
)

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
	stmt, err := Db.Prepare("Insert into users set username = ? , password = ? , firstname = ? , lastname = ? , token = ?")

	if err != nil {
		fmt.Println(err)
		return err
	}

	token, errToken := createToken(user)

	if errToken != nil {
		fmt.Println(err)
		return err
	}

	_, err = stmt.Exec(user.Username, user.Password, user.FirstName, user.LastName, token)

	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

//Check if password is correct and if the password is correct it returns the token corresponding to the user
func CheckPassword(user model.UserLogin) (error, string) {
	query := "Select password from users where username = '" + user.Username + "';"
	var pass string

	err := Db.QueryRow(query).Scan(&pass)

	if err != nil {
		fmt.Println(err)
		return err, ""
	}

	err = bcrypt.CompareHashAndPassword([]byte(pass), []byte(user.Password))

	if err != nil {
		fmt.Println("Bcrypt error ", err)
		return err, ""
	}

	var token string

	err, token = GetToken(user.Username)

	if err != nil {
		fmt.Println(err)
		return err, ""
	}

	err = addTokenToDb(token)

	if err != nil {
		return err, ""
	}

	return nil, token
}

func GetToken(username string) (error, string) {

	query := "select token from users where username = '" + username + "';"

	var token string

	err := Db.QueryRow(query).Scan(&token)

	if err != nil {
		return err, ""
	}
	return nil, token
}

func CheckTokenInDb(token string) (error, bool) {

	query := "select exists(select 1 from user_token where token = '" + token + "');"

	var exist bool

	err := Db.QueryRow(query).Scan(&exist)

	if err != nil {
		fmt.Println(err)
		return err, true
	}

	return nil, exist

}
func createToken(user model.User) (string, error) {
	token, err := bcrypt.GenerateFromPassword([]byte(user.Username+user.Password), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	if err != nil {
		return "", err
	}

	return string(token), nil
}

func addTokenToDb(token string) error {

	stmt, err := Db.Prepare("Insert into user_token set token = ? , startdate = ? ")

	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = stmt.Exec(token, time.Now())

	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("Token added to database")
	return nil
}

func CheckTokenExistsInDb(token string) (bool, error) {
	query := "select exists(select 1 from user_token where token = '" + token + "');"

	var exists bool

	err := Db.QueryRow(query).Scan(&exists)

	if err != nil {
		fmt.Println(err)
		return true, err
	}
	return exists, nil
}

func GetUserDetail(username string) (string, string, error) {

	query := "select firstname, lastname from users where username = '" + username + "';"
	var firstname, lastname string
	err := Db.QueryRow(query).Scan(&firstname, &lastname)

	if err != nil {
		fmt.Println(err)
		return "", "", err
	}
	return firstname, lastname, nil
}
