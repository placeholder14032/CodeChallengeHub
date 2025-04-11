package middleware

import (
	"fmt"
	"net/http"

	"github.com/placeHolder143032/CodeChallengeHub/database"
)

func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	fmt.Print("RequireAuth middleware called\n")
    return func(w http.ResponseWriter, r *http.Request) {
        cookie, err := r.Cookie("session_id")
        if err != nil {
            http.Redirect(w, r, "/login-user", http.StatusSeeOther)
            return
        }

        // Verify session in database
        valid, err := database.ValidateSession(cookie.Value)
        if err != nil || !valid {
            http.Redirect(w, r, "/login-user", http.StatusSeeOther)
            return
        }

        next.ServeHTTP(w, r)
    }
}