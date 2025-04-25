package routes

import (
    "log"
    "net/http"
    "strconv"

    "github.com/placeHolder143032/CodeChallengeHub/database"
    "github.com/placeHolder143032/CodeChallengeHub/middleware"
)

func DraftProblem(w http.ResponseWriter, r *http.Request) {
    // Only allow POST method
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // Check if user is admin
    userIDValue := r.Context().Value(middleware.UserIDKey)
    if userIDValue == nil {
        http.Redirect(w, r, "/login-user", http.StatusSeeOther)
        return
    }

    userID, ok := userIDValue.(int)
    if !ok {
        http.Error(w, "Invalid session", http.StatusInternalServerError)
        return
    }

    // Check admin status
    isAdmin, err := database.GetUserRole(userID)
    if err != nil || isAdmin != 1 {
        http.Error(w, "Unauthorized", http.StatusForbidden)
        return
    }

    // Get problem ID from form
    problemID := r.FormValue("problem_id")
    if problemID == "" {
        http.Error(w, "Problem ID is required", http.StatusBadRequest)
        return
    }

    // Convert problem ID to int
    pid, err := strconv.Atoi(problemID)
    if err != nil {
        http.Error(w, "Invalid problem ID", http.StatusBadRequest)
        return
    }

    // Toggle problem draft status
    err = database.ToggleProblemDraftStatus(pid)
    if err != nil {
        log.Printf("Error toggling draft status: %v", err)
        http.Error(w, "Failed to update problem status", http.StatusInternalServerError)
        return
    }

    // Redirect back to problems list
    http.Redirect(w, r, "/problems", http.StatusSeeOther)
}


// @desc toggle problem publish status
// @route POST /publish-problem
// @access private (admin only)
func PublishProblem(w http.ResponseWriter, r *http.Request) {
    // Only allow POST method
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // Check if user is admin
    userIDValue := r.Context().Value(middleware.UserIDKey)
    if userIDValue == nil {
        http.Redirect(w, r, "/login-admin", http.StatusSeeOther)
        return
    }

    userID, ok := userIDValue.(int)
    if !ok {
        http.Error(w, "Invalid session", http.StatusInternalServerError)
        return
    }

    isAdmin, err := database.GetUserRole(userID)
    if err != nil || isAdmin != 1 {
        http.Error(w, "Unauthorized", http.StatusForbidden)
        return
    }

    // Parse form
    err = r.ParseForm()
    if err != nil {
        http.Error(w, "Unable to parse form", http.StatusBadRequest)
        return
    }

    // Get problem ID from form
    problemID := r.FormValue("problem_id")
    if problemID == "" {
        http.Error(w, "Problem ID is required", http.StatusBadRequest)
        return
    }

    // Convert problem ID to int
    pid, err := strconv.Atoi(problemID)
    if err != nil {
        http.Error(w, "Invalid problem ID", http.StatusBadRequest)
        return
    }

    // Toggle problem publish status
    err = database.PublishProblem(pid)
    if err != nil {
        log.Printf("Error toggling publish status for problem %d: %v", pid, err)
        http.Error(w, "Failed to update problem status", http.StatusInternalServerError)
        return
    }

    // Redirect back to admin problems list
    http.Redirect(w, r, "/problems-admin", http.StatusSeeOther)
}