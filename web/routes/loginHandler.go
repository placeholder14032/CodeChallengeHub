package routes

import (
	"fmt"
	"net/http"

	"github.com/placeHolder143032/CodeChallengeHub/database"
	"github.com/placeHolder143032/CodeChallengeHub/models"
)

// @desc login and check password for signin user
// @route POST /api/auth/login-user
// @access public
func  LoginUser(w http.ResponseWriter, r *http.Request) {
	fmt.Print("Login user")

	// Parse the request body
	username := r.FormValue("username")
	password := r.FormValue("password")
	fmt.Print("Username:", username)
	fmt.Println(", Password:", password)

	// create user
	targetUser := models.User{
		Username: username,
		Password: password,
	}

	id,err := database.SignInUser(targetUser)

	targetUser.ID = id

	// TODO
	if(err!=nil){
		fmt.Println("Error:", err)
	}else{
	// Redirect to profile page with status code 303 (See Other)
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
	}

}

// @desc login and check password for signin admin user
// @route POST /api/auth/login-admin
// @access public
func  LoginAdmin(w http.ResponseWriter, r *http.Request) {
	fmt.Print("Login admin")

	// Parse the request body
	username := r.FormValue("username")
	password := r.FormValue("corp-key")
	fmt.Print("Username:", username)
	fmt.Println(", Password:", password)

	// create user
	targetUser := models.User{
		Username: username,
		Password: password,
	}

	id,err := database.SignInUser(targetUser)

	targetUser.ID = id
	// targetUser.Is_admin = 1

	// TODO
	if(err!=nil){
		fmt.Println("Error:", err)
	}else{
	// Redirect to profile page with status code 303 (See Other)
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
	}
}