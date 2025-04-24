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
    fmt.Print("this is a print in the submit answer function to check if it is working")
    // Parse multipart form
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

    // Save code to file
    codePath, err := saveCodeToFile(code, userID, pid)
    if err != nil {
        log.Printf("SubmitAnswer: Failed to save code: %v", err)
        http.Error(w, "Failed to save code", http.StatusInternalServerError)
        return
    }

    // Create submission record
    submission := models.Submission{
        UserId:    userID,
        ProblemId: pid,
        CodePath:  codePath,
        State:     0, // pending
        CreatedAt: time.Now(),
    }

    // Save to database
    err = database.SubmitCode(submission)
    if err != nil {
        log.Printf("SubmitAnswer: Failed to save submission: %v", err)
        http.Error(w, "Failed to create submission", http.StatusInternalServerError)
        return
    }

    // Redirect to submissions page
    http.Redirect(w, r, "/submissions", http.StatusSeeOther)
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


// @desc get HTML page for a specific submission
// @route GET /submission?id=<submission_id>
// @access private (only accessible to the submission owner or admin)
func GoSubmissionView(w http.ResponseWriter, r *http.Request) {

	// renderTemplate(w, "/submission.html", data)
    renderTemplate(w, "/submission.html", nil)

}