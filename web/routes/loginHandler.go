package routes

import (
	"fmt"
	"log"
	"net/http"

	"github.com/placeHolder143032/CodeChallengeHub/database"
	"github.com/placeHolder143032/CodeChallengeHub/models"
)

// @desc login and check password for signin user
// @route POST /api/auth/login-user
// @access public
func LoginUser(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    username := r.FormValue("username")
    password := r.FormValue("password")

    // Validate input
    if username == "" || password == "" {
        data := struct {
            Error    string
            Username string
        }{
            Error:    "Username and password are required",
            Username: username,
        }
        renderTemplate(w, "auth/userLogin.html", data)
        return
    }

    user := models.User{
        Username: username,
        Password: password,
    }

    // Attempt login
    _, sessionID, err := database.SignInUser(user)
    if err != nil {
        var errorMsg string
        switch err.Error() {
        case "user does not exist":
            errorMsg = "Invalid username or password"
        case "wrong password":
            errorMsg = "Invalid username or password"
        default:
            log.Printf("Login error: %v", err)
            errorMsg = "Error during login"
        }

        data := struct {
            Error    string
            Username string
        }{
            Error:    errorMsg,
            Username: username,
        }
        renderTemplate(w, "auth/userLogin.html", data)
        return
    }

    // Set session cookie
    cookie := &http.Cookie{
        Name:     "session",
        Value:    sessionID,
        Path:     "/",
        HttpOnly: true,
        Secure:   true,
        MaxAge:   3600 * 24, // 24 hours
    }
    http.SetCookie(w, cookie)

    // Redirect to problems page on success
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

	if(err!=nil){
		fmt.Println("Error:", err)
	}else{
	// Redirect to profile page with status code 303 (See Other)
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
	}
}