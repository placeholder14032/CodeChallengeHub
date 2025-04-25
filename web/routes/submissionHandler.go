package routes

import (
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/placeHolder143032/CodeChallengeHub/database"
	"github.com/placeHolder143032/CodeChallengeHub/middleware"
	"github.com/placeHolder143032/CodeChallengeHub/models"

	"fmt"
	"log"
	"os"
)

// @desc submit your answer to a problem
// @route POST /api/submit_answer
// @access private you can only access it if you are logged in
func SubmitAnswer(w http.ResponseWriter, r *http.Request) {
    // fmt.Print("this is a print in the submit answer function to check if it is working")
    err := r.ParseMultipartForm(10 << 20) // 10 MB max memory
    if err != nil {
        log.Printf("SubmitAnswer: Failed to parse form: %v", err)
        http.Error(w, "Failed to parse form", http.StatusBadRequest)
        return
    }

    // Get user ID from context
    userID, ok := r.Context().Value(middleware.UserIDKey).(int)
    if !ok {
        log.Printf("SubmitAnswer: User not authenticated")
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    // Get form values
    problemID := r.FormValue("problem_id")
    if problemID == "" {
        log.Printf("SubmitAnswer: Missing problem ID")
        http.Error(w, "Problem ID is required", http.StatusBadRequest)
        return
    }

    code := r.FormValue("code")
    if code == "" {
        log.Printf("SubmitAnswer: Missing code")
        http.Error(w, "Code is required", http.StatusBadRequest)
        return
    }

    // Convert problemID to int
    pid, err := strconv.Atoi(problemID)
    if err != nil {
        log.Printf("SubmitAnswer: Invalid problem ID: %v", err)
        http.Error(w, "Invalid problem ID", http.StatusBadRequest)
        return
    }

    codePath, err := saveCodeToFile(code, userID, pid)
    if err != nil {
        log.Printf("SubmitAnswer: Failed to save code: %v", err)
        http.Error(w, "Failed to save code", http.StatusInternalServerError)
        return
    }

    submission := models.Submission{
        UserId:    userID,
        ProblemId: pid,
        CodePath:  codePath,
        State:     0, // pending
        CreatedAt: time.Now(),
    }

    err = database.SubmitCode(submission)
    if err != nil {
        log.Printf("SubmitAnswer: Failed to save submission: %v", err)
        http.Error(w, "Failed to create submission", http.StatusInternalServerError)
        return
    }
    http.Redirect(w, r, "/my_submissions", http.StatusSeeOther)
}

func saveCodeToFile(code string, userID, problemID int) (string, error) {
    baseDir := "pkg"
    
    // submissions/user_{userID}/problem_{problemID}
    userDir := fmt.Sprintf("%s/%d/submission", baseDir, userID)
    problemDir := fmt.Sprintf("%s/problem_%d", userDir, problemID)
    
    if err := os.MkdirAll(problemDir, 0755); err != nil {
        return "", fmt.Errorf("failed to create directories: %v", err)
    }
    
    // timestamp for unique name
    timestamp := time.Now().Format("20060102_150405")
    filename := fmt.Sprintf("submission_%s_%d_%d.go", timestamp, userID, problemID)
    filePath := filepath.Join(problemDir, filename)
    
    // Write code to file
    if err := os.WriteFile(filePath, []byte(code), 0644); err != nil {
        return "", fmt.Errorf("failed to write code file: %v", err)
    }
    
    // Verify file was created
    if _, err := os.Stat(filePath); os.IsNotExist(err) {
        return "", fmt.Errorf("file was not created: %v", err)
    }
    
    absPath, err := filepath.Abs(filePath)
    if err != nil {
        return "", fmt.Errorf("failed to get absolute path: %v", err)
    }
    
    return absPath, nil
}


// @desc get Html page for submitting problem
// @route GET / submit_answer
// @access private (you can only access this page if you are logged in)
func GoSubmitAnswer(w http.ResponseWriter, r *http.Request) {
    problemID := r.URL.Query().Get("problem")
    if problemID == "" {
        http.Error(w, "Problem ID is required", http.StatusBadRequest)
        return
    }

    data := struct {
        ProblemID string
    }{
        ProblemID: problemID,
    }

    renderTemplate(w, "problem_submit.html", data)
}

// @desc get HTML page for viewing user submissions
// @route GET /my_submissions
// @access private (only accessible to the logged-in user)
func ViewSubmissionsByUser(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userIDValue := r.Context().Value(middleware.UserIDKey)
	if userIDValue == nil {
		log.Printf("ViewSubmissionsByUser: No user ID in context")
		http.Redirect(w, r, "/login-user", http.StatusSeeOther)
		return
	}

	userID, ok := userIDValue.(int)
	if !ok {
		log.Printf("ViewSubmissionsByUser: Invalid user ID type in context")
		http.Error(w, "Invalid session", http.StatusInternalServerError)
		return
	}

	// Fetch all submissions for the user
	submissions, err := database.GetAllSubmissionsByUser(userID)
	if err != nil {
		log.Printf("ViewSubmissionsByUser: Error fetching submissions: %v", err)
		http.Error(w, "Failed to fetch submissions", http.StatusInternalServerError)
		return
	}

	// Get page number from query parameter
	pageStr := r.URL.Query().Get("page")
	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	// Define items per page
	const itemsPerPage = 10

	// Calculate pagination details
	totalItems := len(submissions)
	totalPages := (totalItems + itemsPerPage - 1) / itemsPerPage

	// Ensure the current page is within bounds
	if page > totalPages {
		page = totalPages
	}
	if page < 1 {
		page = 1
	}

	// Calculate start and end indices for the current page
	start := (page - 1) * itemsPerPage
	end := start + itemsPerPage
	if end > totalItems {
		end = totalItems
	}

	// Get the submissions for the current page
	paginatedSubmissions := submissions[start:end]

	// Generate page numbers for navigation
	var pageNumbers []int
	for i := 1; i <= totalPages; i++ {
		pageNumbers = append(pageNumbers, i)
	}

	// Check if user is admin
	isAdmin, err := database.GetUserRole(userID)
	if err != nil {
		log.Printf("ViewSubmissionsByUser: Error checking user role: %v", err)
		http.Error(w, "Failed to check user role", http.StatusInternalServerError)
		return
	}

    // Convert states to status strings for the paginated submissions
    for i, sub := range paginatedSubmissions {
        paginatedSubmissions[i].Status = database.GetStatusFromState(sub.State)
    }

    // Prepare template data
    data := struct {
        Submissions  []models.Submission
        CurrentPage  int
        PrevPage     int
        NextPage     int
        TotalPages   int
        PageNumbers  []int
        IsAdmin      bool
        CurrentUser  int
    }{
        Submissions:  paginatedSubmissions,  // Use paginated submissions
        CurrentPage:  page,
        PrevPage:     page - 1,
        NextPage:     page + 1,
        TotalPages:   totalPages,
        PageNumbers:  pageNumbers,
        IsAdmin:      isAdmin == 1,
        CurrentUser:  userID,
    }

    // Use the correct template path without leading slash
    renderTemplate(w, "mySubmissions.html", data)
}

// @desc get HTML page for a specific submission
// @route GET /submission?id=<submission_id>
// @access private (only accessible to the submission owner or admin)
func GoSubmissionView(w http.ResponseWriter, r *http.Request) {
	// renderTemplate(w, "/submission.html", data)
    renderTemplate(w, "/submission.html", nil)

}