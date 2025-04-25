package routes

import (
	"net/http"
	"strconv"
	"github.com/placeHolder143032/CodeChallengeHub/database"
	"github.com/placeHolder143032/CodeChallengeHub/middleware"
)

// @desc renders the add admin page.
// @route GET /add-admin
// @access private, admin only
func GoAddAdminPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "addAdmin.html", nil)

}

// @desc renders the remove admin page.
// @route GET /remove-admin
// @access private, admin only
func GoRemoveAdminPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "removeAdmin.html", nil)
}

// MakeAdmin handles the form submission to add admin role.
func MakeAdmin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userIDStr := r.FormValue("userID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	err = database.MakeAdmin(userID)
	if err != nil {
		if err.Error() == "user does not exist" {
			http.Error(w, "User does not exist", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to add admin role", http.StatusInternalServerError)
		}
		return
	}

	// Redirect to profile page or display success message
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

// RemoveAdmin handles the form submission to remove admin role.
func RemoveAdmin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userIDStr := r.FormValue("userID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Prevent self-demotion (optional, depending on requirements)
	ctx := r.Context()
	currentUserID, ok := ctx.Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if currentUserID == userID {
		http.Error(w, "Cannot remove your own admin role", http.StatusForbidden)
		return
	}

	err = database.RemoveAdmin(userID)
	if err != nil {
		if err.Error() == "user does not exist" {
			http.Error(w, "User does not exist", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to remove admin role", http.StatusInternalServerError)
		}
		return
	}

	// Redirect to profile page or display success message
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}