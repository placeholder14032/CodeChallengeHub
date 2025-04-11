package routes

import (
	"fmt"
	"net/http"

	"github.com/placeHolder143032/CodeChallengeHub/database"
	"github.com/placeHolder143032/CodeChallengeHub/models"

	"golang.org/x/crypto/bcrypt"
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

		hashedPassword, err := HashPassword(password)
		if err != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			return
		}

		// fmt.Println("Hashed Password:", hashedPassword)

		// create user
		createdUser := models.User{
			Username: username,
			Password: hashedPassword,
		}

		database.SignUpUser(createdUser)
		
	// Redirect to profile page with status code 303 (See Other)
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
	} else{
		passwordMismatch()
	}
}

func HashPassword(password string) (string, error) {
    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    return string(hash), nil
}

// TODO
func passwordMismatch(){}