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
    fmt.Print("Login user hereee 1\n")
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

    fmt.Printf("Usernameeeeeeeee: %s, Password: %s\n", username   , password)

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

    fmt.Printf("Session IDDDDDDDDDD: %s\n", sessionID)
    // Set session cookie
    cookie := &http.Cookie{
        Name:     "session_id",
        Value:    sessionID,
        Path:     "/",
        HttpOnly: true,
        Secure:   true,
        MaxAge:   3600 * 24, // 24 hours
        SameSite: http.SameSiteLaxMode,
    }
    fmt.Printf("Cookieeeeeeee: %v\n", cookie)
    http.SetCookie(w, cookie)
    fmt.Printf("Cookie seeeeeeeeet: %s\n", cookie.String())

    // Redirect to problems page on success
    http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

// @desc login and check password for signin admin user
// @route POST /api/auth/login-admin
// @access public
func LoginAdmin(w http.ResponseWriter, r *http.Request) {
    fmt.Print("Login admin\n")

    // Parse the request body
    username := r.FormValue("username")
    password := r.FormValue("corp-key")
    fmt.Printf("Username: %s, Password: %s\n", username, password)

    targetUser := models.User{
        Username: username,
        Password: password,
    }

    id, sessionID, err := database.SignInUser(targetUser)
    fmt.Printf("Session ID: %s\n", sessionID)

    targetUser.ID = id

    if err != nil {
        fmt.Printf("Error: %v\n", err)
        data := struct {
            Error    string
            Username string
        }{
            Error:    "Invalid username or password",
            Username: username,
        }
        renderTemplate(w, "auth/adminLogin.html", data) 
        return
    }

    // Set session cookie
    cookie := &http.Cookie{
        Name:     "session_id",
        Value:    sessionID,
        Path:     "/",
        HttpOnly: true,
        Secure:   true,
        MaxAge:   3600 * 24, // 24 hours
        SameSite: http.SameSiteLaxMode,
    }
    http.SetCookie(w, cookie)

    http.Redirect(w, r, "/profile", http.StatusSeeOther)
}