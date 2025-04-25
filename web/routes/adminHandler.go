package routes

import (
    "log"
    "net/http"

    "github.com/placeHolder143032/CodeChallengeHub/database"
    "github.com/placeHolder143032/CodeChallengeHub/middleware"
	"github.com/placeHolder143032/CodeChallengeHub/models"
)

func GoAdminControls(w http.ResponseWriter, r *http.Request) {
    // Get current user ID from context
    userID, ok := r.Context().Value(middleware.UserIDKey).(int)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    // Verify admin status
    isAdmin, err := database.GetUserRole(userID)
    if err != nil || isAdmin != 1 {
        log.Printf("Non-admin user %d attempted to access admin controls", userID)
        http.Redirect(w, r, "/profile", http.StatusSeeOther)
        return
    }

    // Get all users
    users, err := database.GetAllUsers()
    if err != nil {
        log.Printf("Error fetching users: %v", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    // Prepare template data
    data := struct {
        Users    []models.User
        IsAdmin  bool
        UserID   int
    }{
        Users:    users,
        IsAdmin:  true,
        UserID:   userID,
    }

    // Render template
    renderTemplate(w, "admin_controls.html", data)
}