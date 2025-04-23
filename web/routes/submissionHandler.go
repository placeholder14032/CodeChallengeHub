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
        UserId:    userID,
        ProblemId: req.ProblemID,
        CodePath:  codePath,
        State:      0, // 0 could represent "pending" state
        CreatedAt: time.Now(),
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


// @desc get HTML page for a specific submission
// @route GET /submission?id=<submission_id>
// @access private (only accessible to the submission owner or admin)
func GoSubmissionView(w http.ResponseWriter, r *http.Request) {
	// Assume user is authenticated
	currentUserID := "user123" // TODO: Replace with actual user ID from session/token

	// Get submission ID from query parameter
	submissionID := r.URL.Query().Get("id")
	if submissionID == "" {
		http.Error(w, "Submission ID is required", http.StatusBadRequest)
		return
	}

	// TODO: Fetch submission from database
	// For now, using static data
	submissions := []Submission{
		{
			ID:            "s1",
			ProblemID:     "1",
			ProblemTitle:  "Two Sum",
			UserID:        "user123",
			OwnerUsername: "john_doe",
			Code:          "func twoSum(nums []int, target int) []int {\n    for i := 0; i < len(nums); i++ {\n        for j := i + 1; j < len(nums); j++ {\n            if nums[i] + nums[j] == target {\n                return []int{i, j}\n            }\n        }\n    }\n    return nil\n}",
			Status:        "OK",
			TimeUsed:      50, // Example time in ms
			SubmittedAt:   time.Now().Add(-24 * time.Hour),
		},
		{
			ID:            "s2",
			ProblemID:     "2",
			ProblemTitle:  "Add Two Numbers",
			UserID:        "user456",
			OwnerUsername: "jane_smith",
			Code:          "func addTwoNumbers(l1 *ListNode, l2 *ListNode) *ListNode {\n    // ...",
			Status:        "Wrong Answer",
			TimeUsed:      120,
			SubmittedAt:   time.Now().Add(-20 * time.Hour),
		},
	}

	// Find the submission
	var submission Submission
	found := false
	for _, s := range submissions {
		if s.ID == submissionID {
			submission = s
			found = true
			break
		}
	}

	if !found {
		http.Error(w, "Submission not found", http.StatusNotFound)
		return
	}

	// Access control: only the owner or an admin can view
	isAdmin := false // TODO: Implement admin check (e.g., from user role in DB)
	if submission.UserID != currentUserID && !isAdmin {
		http.Error(w, "Forbidden: You can only view your own submissions", http.StatusForbidden)
		return
	}

	// Prepare data for template
	data := struct {
		Submission Submission
	}{
		Submission: submission,
	}

	renderTemplate(w, "/submission.html", data)
}