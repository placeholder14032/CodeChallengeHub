package routes

import (
	"net/http"

	"github.com/placeHolder143032/CodeChallengeHub/database"
	"github.com/placeHolder143032/CodeChallengeHub/models"
)

// @desc get Html page for profile page
// @route GET /profile
// @access private (you can only access this page if you are logged in)
func GoProfilePage(w http.ResponseWriter, r *http.Request) {
    // Get session cookie
    cookie, err := r.Cookie("session_id")
    if err != nil {
        http.Redirect(w, r, "/login-user", http.StatusSeeOther)
        return
    }

    // Get current user
    user, err := database.GetCurrentUser(cookie.Value)
    if err != nil {
        http.Error(w, "Failed to get user profile", http.StatusInternalServerError)
        return
    }

    // Calculate success rate
    var successRate float64
    if user.AttemptedProblems > 0 {
        successRate = float64(user.SolvedProblems) / float64(user.AttemptedProblems) * 100
    }

    // Prepare data for template
    data := models.User{
        Username:        user.Username,
        SuccessRate:     successRate,
        SolvedProblems: user.SolvedProblems,
        IsAdmin:         user.IsAdmin,
        ID:          user.ID,
    }

    renderTemplate(w, "profilePage.html", data)
}