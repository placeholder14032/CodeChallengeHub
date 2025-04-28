package middleware

import (
	"net/http"

	"github.com/placeHolder143032/CodeChallengeHub/database"
    "context"

	"log"
    "fmt"
)
type contextKey string

const UserIDKey contextKey = "userID"

func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        cookie, err := r.Cookie("session_id") 
        fmt.Println("Cookie: ", cookie)
        fmt.Println("ERRR ", err)
        if err != nil {
            log.Printf("No session cookie for request %s: %v", r.URL.String(), err)
            http.Redirect(w, r, "/login-user", http.StatusSeeOther)
            return
        }

        log.Printf("Found cookie: session_id=%s for request %s", cookie.Value, r.URL.String())
        userID, valid, err := database.ValidateSession(cookie.Value)
        if err != nil {
            log.Printf("Session validation error for session %s on request %s: %v", 
                cookie.Value, r.URL.String(), err)
            http.Redirect(w, r, "/login-user", http.StatusSeeOther)
            return
        }
        if !valid {
            log.Printf("Invalid or expired session %s for request %s", cookie.Value, r.URL.String())
            http.Redirect(w, r, "/login-user", http.StatusSeeOther)
            return
        }

        log.Printf("Session validated: userID=%d for session %s on request %s", 
            userID, cookie.Value, r.URL.String())
        ctx := context.WithValue(r.Context(), UserIDKey, userID)
        next.ServeHTTP(w, r.WithContext(ctx))
    }
}