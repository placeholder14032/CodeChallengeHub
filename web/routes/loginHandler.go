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
}

// @desc login and check password for signin admin user
// @route POST /api/auth/login-admin
// @access public
func  LoginAdmin(w http.ResponseWriter, r *http.Request) {
	fmt.Print("Loginnnnnnn Adminnnnn")
}