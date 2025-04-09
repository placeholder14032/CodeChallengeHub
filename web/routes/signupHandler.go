package routes

import (
	"fmt"
	"net/http"

	"github.com/placeHolder143032/CodeChallengeHub/database"
	"github.com/placeHolder143032/CodeChallengeHub/models"
)

// @desc creaing a new user account for signup
// @route POST /api/auth/register-user
// @access public
func  SignupUser(w http.ResponseWriter, r *http.Request) {
	fmt.Print("Signup user:")

	username := r.FormValue("username")
	password := r.FormValue("password")
	passwordConfirm := r.FormValue("passwordConfirm")

	fmt.Print("Username:", username)
	fmt.Print(", Password:", password)
	fmt.Println(", Confirm Password:", passwordConfirm)

	if (password==passwordConfirm){
		// ok continue signing up

		// create user
		createdUser := models.User{
			Username: username,
			Password: password,
		}

		database.SignUpUser(createdUser)
		
	// Redirect to profile page with status code 303 (See Other)
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
	} else{
		passwordMismatch()
	}
}

// TODO
func passwordMismatch(){}