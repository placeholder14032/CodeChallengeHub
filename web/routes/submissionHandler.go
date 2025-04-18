package routes

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"time"

	"github.com/placeHolder143032/CodeChallengeHub/database"
	"github.com/placeHolder143032/CodeChallengeHub/middleware"
	"github.com/placeHolder143032/CodeChallengeHub/models"

	"fmt"
	"os"
)

// we should make it compatible and same as pooya's
type Submission struct {
    ID             string
    ProblemID      string
    ProblemTitle   string 
    UserID         string
    OwnerUsername  string 
    Code           string
    Status         string // "Pending", "OK", "Compile Error", "Wrong Answer", "Memory Limit", "Time Limit", "Runtime Error"
    TimeUsed       int  
    SubmittedAt    time.Time
}

// @desc submit your answer to a problem
// @route POST /api/submit_answer
// @access private you can only access it if you are logged in
func SubmitAnswer(w http.ResponseWriter, r *http.Request) {
    fmt.Print("wsdfgcjhcv kghc utdltuodsutkgfhgftckjtuxy")
    // Get user ID from context (set by RequireAuth middleware)
    userID, ok := r.Context().Value(middleware.UserIDKey).(int)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    // Parse request body
    type submitRequest struct {
        ProblemID int    `json:"problem_id"`
        Code      string `json:"code"`
    }

    var req submitRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Validate input
    if req.ProblemID <= 0 || req.Code == "" {
        http.Error(w, "Missing required fields", http.StatusBadRequest)
        return
    }

    // Generate code file path (you might want to implement your own file storage logic)
    codePath, err := saveCodeToFile(req.Code, userID, req.ProblemID)
    if err != nil {
        http.Error(w, "Failed to save code", http.StatusInternalServerError)
        return
    }

    // Create submission record
    submission := models.Submission{
        User_id:    userID,
        Problem_id: req.ProblemID,
        Code_path:  codePath,
        State:      0, // 0 could represent "pending" state
        Created_at: time.Now(),
        // Runtime_ms, Memory_used, and Error_message will be updated after evaluation
    }

    // Save submission to database
    err = database.SubmitCode(submission)
    if err != nil {
        http.Error(w, "Failed to create submission", http.StatusInternalServerError)
        return
    }

    // submition evaluation

    // Send response
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "message":     "Submission received",
        "submission":  submission,
    })
}

func saveCodeToFile(code string, userID, problemID int) (string, error) {
    // Define base directory for code storage
    baseDir := "pkg"
    
    // Create directory structure: submissions/user_{userID}/problem_{problemID}
    userDir := fmt.Sprintf("%s/%d/submission", baseDir, userID)
    problemDir := fmt.Sprintf("%s/problem_%d", userDir, problemID)
    
    // Ensure directories exist
    if err := os.MkdirAll(problemDir, 0755); err != nil {
        return "", fmt.Errorf("failed to create directories: %v", err)
    }
    
    // Generate unique filename using timestamp
    timestamp := time.Now().Format("20060102_150405") // Format: YYYYMMDD_HHMMSS
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
    
    // Return absolute path
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
	renderTemplate(w, "problem_submit.html", nil)
}