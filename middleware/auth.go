package middleware

import (
	"net/http"

	"context"

	"github.com/placeHolder143032/CodeChallengeHub/database"

	"log"
)
type contextKey string

const UserIDKey contextKey = "userID"

func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        cookie, err := r.Cookie("session")
        if err != nil {
            log.Printf("No session cookie found: %v", err)
            http.Redirect(w, r, "/login-user", http.StatusSeeOther)
            return
        }

        userID, err := database.GetUserIDFromSession(cookie.Value)
        if err != nil {
            log.Printf("Failed to get user ID from session: %v", err)
            http.Redirect(w, r, "/login-user", http.StatusSeeOther)
            return
        }

        log.Printf("Session %s validated, user_id=%d", cookie.Value, userID)
        
        ctx := context.WithValue(r.Context(), UserIDKey, userID)
        next.ServeHTTP(w, r.WithContext(ctx))
    }
}