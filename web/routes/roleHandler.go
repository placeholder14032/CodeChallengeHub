package routes

import (
    "log"
    "net/http"
    "strconv"

    "github.com/placeHolder143032/CodeChallengeHub/database"
    "github.com/placeHolder143032/CodeChallengeHub/middleware"
)

func ToggleAdminRole(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // Get current user ID from context and verify admin status
    currentUserID, ok := r.Context().Value(middleware.UserIDKey).(int)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    // Verify current user is admin
    isAdmin, err := database.GetUserRole(currentUserID)
    if err != nil || isAdmin != 1 {
        log.Printf("Unauthorized admin role change attempt by user %d", currentUserID)
        http.Error(w, "Only administrators can modify user roles", http.StatusForbidden)
        return
    }

    // Get target user ID from form
    targetUserID := r.FormValue("userID")
    if targetUserID == "" {
        http.Error(w, "User ID is required", http.StatusBadRequest)
        return
    }

    // Convert target user ID to int
    uid, err := strconv.Atoi(targetUserID)
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    // Change user role
    err = database.ChangeUserRole(uid)
    if err != nil {
        log.Printf("Error changing user role: %v", err)
        http.Error(w, "Failed to change user role", http.StatusInternalServerError)
        return
    }

    // Redirect back to profile page
    http.Redirect(w, r, "/profile?id="+targetUserID, http.StatusSeeOther)
}