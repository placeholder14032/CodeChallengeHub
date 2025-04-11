package routes

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/placeHolder143032/CodeChallengeHub/database"
	"github.com/placeHolder143032/CodeChallengeHub/models"
)

// @desc login and check password for signin user
// @route POST /api/auth/login-user
// @access public
func LoginUser(w http.ResponseWriter, r *http.Request) {
	// fmt.Print("alskvhlafshvha")
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    username := r.FormValue("username")
    password := r.FormValue("password")

    user := models.User{
        Username: username,
        Password: password,
    }

    userID, sessionID, err := database.SignInUser(user)
    if err != nil {
        // http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		errors.New(err.Error())
		fmt.Print(err)
        return
    }

	fmt.Println("User ID:", userID) //  idk how to use it or write it in db

    // Set session cookie
    cookie := http.Cookie{
        Name:     "session_id",
        Value:    sessionID,
        Path:     "/",
        HttpOnly: true,
        Secure:   false, // Set to true in production later
        MaxAge:   24 * 60 * 60, // 1 day
    }

    http.SetCookie(w, &cookie)

    // Redirect to profile page
    http.Redirect(w, r, "/profile", http.StatusSeeOther)
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

	id,sessionID,err := database.SignInUser(targetUser)

	fmt.Print("Session ID:", sessionID) //  idk how to use it for now

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