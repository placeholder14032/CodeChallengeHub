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
	username := r.FormValue("username")
	password := r.FormValue("password")
	passwordConfirm := r.FormValue("passwordConfirm")
	fmt.Println("Username:", username)
	fmt.Println("Password:", password)
	fmt.Println("Confirm Password:", passwordConfirm)
}


// @desc creaing a new admin account for signup ????????
// @route POST /api/auth/register-user
// @access public
func  SignupAdmin(w http.ResponseWriter, r *http.Request) {
	fmt.Print("Signinnnnnnnn")
	username := r.FormValue("username")
	corpEmail := r.FormValue("corp-email")
	corpKey := r.FormValue("corp-key")
	fmt.Println("Username:", username)

	fmt.Println("Corp Email:", corpEmail)
	fmt.Println("Corp Key:", corpKey)

}

func checkUserStatus_Signup(username, password string) bool {	
	// check id user exists in the database
	// check if the password is same as the confirm password
	return true
}
