package model

type User struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

type UserLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserProfile struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}
