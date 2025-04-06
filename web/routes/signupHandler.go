package routes

import(
	"fmt"
	"net/http"
)

// @desc creaing a new user account for signup
// @route POST /api/auth/register-user
// @access public
func  SignupUser(w http.ResponseWriter, r *http.Request) {
	fmt.Print("Signinnnnnnnn")
}


// @desc creaing a new admin account for signup ????????
// @route POST /api/auth/register-user
// @access public
func  SignupAdmin(w http.ResponseWriter, r *http.Request) {
	fmt.Print("Signinnnnnnnn")
}