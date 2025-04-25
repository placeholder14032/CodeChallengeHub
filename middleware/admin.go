package middleware

import (
	"context"
	"net/http"
	"github.com/placeHolder143032/CodeChallengeHub/database"
	"log"
)

// RequireAdmin ensures the user is authenticated and has admin privileges.
func RequireAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// First, check if the user is authenticated
		cookie, err := r.Cookie("session_id")
		if err != nil {
			log.Printf("No session cookie for request %s: %v", r.URL.String(), err)
			http.Redirect(w, r, "/login-admin", http.StatusSeeOther)
			return
		}

		userID, valid, err := database.ValidateSession(cookie.Value)
		if err != nil {
			log.Printf("Session validation error for session %s on request %s: %v", 
				cookie.Value, r.URL.String(), err)
			http.Redirect(w, r, "/login-admin", http.StatusSeeOther)
			return
		}
		if !valid {
			log.Printf("Invalid or expired session %s for request %s", cookie.Value, r.URL.String())
			http.Redirect(w, r, "/login-admin", http.StatusSeeOther)
			return
		}

		// Check if the user is an admin
		isAdmin, err := database.GetUserRole(userID)
		if err != nil {
			log.Printf("Error checking admin status for user %d: %v", userID, err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		if isAdmin != 1 {
			log.Printf("User %d is not an admin for request %s", userID, r.URL.String())
			http.Error(w, "Forbidden: Admin access required", http.StatusForbidden)
			return
		}

		// Add userID to context and proceed
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}