package middleware

import (
	"net/http"

	"github.com/placeHolder143032/CodeChallengeHub/database"
    "context"
)
type contextKey string

const UserIDKey contextKey = "userID"

func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			http.Redirect(w, r, "/login-user", http.StatusSeeOther)
			return
		}

		userID, valid, err := database.ValidateSession(cookie.Value)
		if err != nil || !valid {
			http.Redirect(w, r, "/login-user", http.StatusSeeOther)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}