package routes

import (
	"log"
	"net/http"

	"github.com/placeHolder143032/CodeChallengeHub/database"
	"github.com/placeHolder143032/CodeChallengeHub/models"

	"golang.org/x/crypto/bcrypt"
)

// @desc creaing a new user account for signup
// @route POST /api/auth/register-user
// @access public
func SignupUser(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    username := r.FormValue("username")
    password := r.FormValue("password")
    passwordConfirm := r.FormValue("passwordConfirm")

    // Validation
    if len(username) < 3 {
        renderTemplate(w, "auth/userSignup.html", struct {
            Error    string
            Username string
        }{
            Error:    "Username must be at least 3 characters long",
            Username: username,
        })
        return
    }

    if len(password) < 8 {
        renderTemplate(w, "auth/userSignup.html", struct {
            Error    string
            Username string
        }{
            Error:    "Password must be at least 8 characters long",
            Username: username,
        })
        return
    }

    if password != passwordConfirm {
        renderTemplate(w, "auth/userSignup.html", struct {
            Error    string
            Username string
        }{
            Error:    "Passwords do not match",
            Username: username,
        })
        return
    }

    // Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        log.Printf("Error hashing password: %v", err)
        renderTemplate(w, "auth/userSignup.html", struct {
            Error    string
            Username string
        }{
            Error:    "Internal server error",
            Username: username,
        })
        return
    }

    user := models.User{
        Username: username,
        Password: string(hashedPassword),
    }

    // Save user to database
    err = database.SignUpUser(user)
    if err != nil {
        var errorMsg string
        if err.Error() == "user already exists" {
            errorMsg = "Username already taken"
        } else {
            errorMsg = "Error creating account"
            log.Printf("Error creating user: %v", err)
        }

        renderTemplate(w, "auth/userSignup.html", struct {
            Error    string
            Username string
        }{
            Error:    errorMsg,
            Username: username,
        })
        return
    }

    // Redirect to login page on success
    http.Redirect(w, r, "/login-user", http.StatusSeeOther)
}



