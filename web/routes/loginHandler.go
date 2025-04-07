package routes

import (
	"fmt"
	"net/http"
)

// @desc login and check password for signin user
// @route POST /api/auth/login-user
// @access public
func  LoginUser(w http.ResponseWriter, r *http.Request) {
	fmt.Print("Loginnnnnnn Userrrrrr")
	// Parse the request body
	username := r.FormValue("username")
	password := r.FormValue("password")
	fmt.Println("Username:", username)
	fmt.Println("Password:", password)
}

// @desc login and check password for signin admin user
// @route POST /api/auth/login-admin
// @access public
func  LoginAdmin(w http.ResponseWriter, r *http.Request) {
	fmt.Print("Loginnnnnnn Adminnnnn")
	username := r.FormValue("username")
	password := r.FormValue("password")
	fmt.Println("Username:", username)
	fmt.Println("Password:", password)
}


func  checkUserStatus_Login(username, password string) bool {
	// check id user exists in the database
	// check if the password is correct
	return true
}