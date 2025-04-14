package routes

import (
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/placeHolder143032/CodeChallengeHub/database"
	"github.com/placeHolder143032/CodeChallengeHub/middleware"
	"github.com/placeHolder143032/CodeChallengeHub/models"
)

// FormData holds form values for repopulation on error
type FormData struct {
	Title       string
	Statement   string
	TimeLimit   string
	MemoryLimit string
	Input       string
	Output      string
}

// @desc add problem to the database
// @route GET,POST /add_problem
// @access private
func AddProblem(w http.ResponseWriter, r *http.Request) {
	log.Printf("AddProblem: Request received, Method: %s, URL: %s", r.Method, r.URL.Path)

	// Parse the HTML template
	tmpl, err := template.ParseFiles("ui/html/add_problem.html")
	if err != nil {
		log.Printf("AddProblem: Template parse error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	log.Printf("AddProblem: Template parsed successfully")

	// Handle GET request (render empty form)
	if r.Method == http.MethodGet {
		log.Printf("AddProblem: Handling GET request")
		err := tmpl.Execute(w, nil)
		if err != nil {
			log.Printf("AddProblem: Template execute error: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		} else {
			log.Printf("AddProblem: GET request served, form rendered")
		}
		return
	}

	// Handle POST request
	if r.Method != http.MethodPost {
		log.Printf("AddProblem: Invalid method %s, expected POST", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	log.Printf("AddProblem: Handling POST request")

	// Parse form data
	err = r.ParseForm()
	if err != nil {
		log.Printf("AddProblem: Form parse error: %v", err)
		tmpl.Execute(w, struct{ Error string }{Error: "Unable to parse form"})
		return
	}
	log.Printf("AddProblem: Form parsed successfully")

	// Extract form values
	form := FormData{
		Title:       r.FormValue("title"),
		Statement:   r.FormValue("statement"),
		TimeLimit:   r.FormValue("time_limit"),
		MemoryLimit: r.FormValue("memory_limit"),
		Input:       r.FormValue("input"),
		Output:      r.FormValue("output"),
	}
	log.Printf("AddProblem: Form data extracted: Title=%s, TimeLimit=%s, MemoryLimit=%s",
		form.Title, form.TimeLimit, form.MemoryLimit)

	// Validate required fields
	if form.Title == "" || form.Statement == "" ||
		form.TimeLimit == "" || form.MemoryLimit == "" || form.Input == "" || form.Output == "" {
		log.Printf("AddProblem: Validation failed: Missing required fields")
		tmpl.Execute(w, struct{ Error string; Form FormData }{
			Error: "All fields are required",
			Form:  form,
		})
		return
	}
	log.Printf("AddProblem: Required fields validated")

	// Validate time limit
	timeLimit, err := strconv.Atoi(form.TimeLimit)
	if err != nil || timeLimit <= 0 {
		log.Printf("AddProblem: Validation failed: Invalid time limit: %s, error: %v", form.TimeLimit, err)
		tmpl.Execute(w, struct{ Error string; Form FormData }{
			Error: "Time limit must be a positive number",
			Form:  form,
		})
		return
	}
	log.Printf("AddProblem: Time limit validated: %d ms", timeLimit)

	// Validate memory limit
	memoryLimit, err := strconv.Atoi(form.MemoryLimit)
	if err != nil || memoryLimit <= 0 {
		log.Printf("AddProblem: Validation failed: Invalid memory limit: %s, error: %v", form.MemoryLimit, err)
		tmpl.Execute(w, struct{ Error string; Form FormData }{
			Error: "Memory limit must be a positive number",
			Form:  form,
		})
		return
	}
	log.Printf("AddProblem: Memory limit validated: %d MB", memoryLimit)

	// Get user ID from context
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		log.Printf("AddProblem: User ID not found in context")
		tmpl.Execute(w, struct{ Error string; Form FormData }{
			Error: "Authentication error",
			Form:  form,
		})
		return
	}
	log.Printf("AddProblem: User ID retrieved: %d", userID)

	// Create directory for problem files
	basePath := "problems"
	year := time.Now().Year()
	userDir := fmt.Sprintf("user_%d", userID)
	timestamp := time.Now().Unix()
	problemDir := filepath.Join(basePath, fmt.Sprintf("%d", year), userDir, fmt.Sprintf("problem_%d_%d", userID, timestamp))

	log.Printf("AddProblem: Creating directory: %s", problemDir)
	if err := os.MkdirAll(problemDir, 0755); err != nil {
		log.Printf("AddProblem: Failed to create directory %s: %v", problemDir, err)
		tmpl.Execute(w, struct{ Error string; Form FormData }{
			Error: "Failed to save problem files",
			Form:  form,
		})
		return
	}
	log.Printf("AddProblem: Directory created successfully")

	// Define file paths
	descriptionPath := filepath.Join(problemDir, "description.txt")
	inputPath := filepath.Join(problemDir, "input.txt")
	outputPath := filepath.Join(problemDir, "output.txt")
	log.Printf("AddProblem: File paths defined: description=%s, input=%s, output=%s",
		descriptionPath, inputPath, outputPath)

	// Save problem content to files
	log.Printf("AddProblem: Writing description file")
	if err := os.WriteFile(descriptionPath, []byte(form.Statement), 0644); err != nil {
		log.Printf("AddProblem: Failed to write description file %s: %v", descriptionPath, err)
		tmpl.Execute(w, struct{ Error string; Form FormData }{
			Error: "Failed to save problem files",
			Form:  form,
		})
		return
	}
	log.Printf("AddProblem: Description file written")

	log.Printf("AddProblem: Writing input file")
	if err := os.WriteFile(inputPath, []byte(form.Input), 0644); err != nil {
		log.Printf("AddProblem: Failed to write input file %s: %v", inputPath, err)
		tmpl.Execute(w, struct{ Error string; Form FormData }{
			Error: "Failed to save problem files",
			Form:  form,
		})
		return
	}
	log.Printf("AddProblem: Input file written")

	log.Printf("AddProblem: Writing output file")
	if err := os.WriteFile(outputPath, []byte(form.Output), 0644); err != nil {
		log.Printf("AddProblem: Failed to write output file %s: %v", outputPath, err)
		tmpl.Execute(w, struct{ Error string; Form FormData }{
			Error: "Failed to save problem files",
			Form:  form,
		})
		return
	}
	log.Printf("AddProblem: Output file written")

	// Create problem struct
	problem := models.Problem{
		UserID:          userID,
		Title:           form.Title,
		DescriptionPath: descriptionPath,
		InputPath:       inputPath,
		OutputPath:      outputPath,
		CreatedTime:     time.Now(),
		IsPublished:     false,
		TimeLimit:       timeLimit,
		MemoryLimit:     memoryLimit,
	}
	log.Printf("AddProblem: Problem struct created: Title=%s, UserID=%d", problem.Title, problem.UserID)

	// Save problem to database
	log.Printf("AddProblem: Saving problem to database")
	if err := database.AddProblem(userID, problem); err != nil {
		log.Printf("AddProblem: Failed to save problem to database: %v", err)
		tmpl.Execute(w, struct{ Error string; Form FormData }{
			Error: "Failed to save problem",
			Form:  form,
		})
		// Cleanup files on database failure
		log.Printf("AddProblem: Cleaning up directory %s due to database failure", problemDir)
		os.RemoveAll(problemDir)
		return
	}
	log.Printf("AddProblem: Problem saved to database")

	// Log success
	log.Printf("AddProblem: Problem '%s' created successfully by user %d", form.Title, userID)

	// Redirect to problems list
	log.Printf("AddProblem: Redirecting to /problems")
	// http.Redirect(w, r, "/problems", http.StatusSeeOther)
}

// @desc get HTML page for all problems with pagination
// @route GET /problems?page=<number>
// @access private (you can only access this page if you are logged in)
func GoProblemsListPageUser(w http.ResponseWriter, r *http.Request) {
	// Get page number from query parameters
	pageStr := r.URL.Query().Get("page")
	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	itemsPerPage := 10
	problems, err := database.GetProblemsPageUser(page, itemsPerPage)
	if err != nil {
		log.Printf("Error fetching problems: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Get total count for pagination
	totalCount, err := database.GetTotalProblemsCount()
	if err != nil {
		log.Printf("Error getting total count: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	totalPages := (totalCount + itemsPerPage - 1) / itemsPerPage

	data := struct {
		Problems    []models.Problem
		CurrentPage int
		PrevPage    int
		NextPage    int
		TotalPages  int
		PageNumbers []int
	}{
		Problems:    problems,
		CurrentPage: page,
		PrevPage:    page - 1,
		NextPage:    page + 1,
		TotalPages:  totalPages,
		PageNumbers: generatePageNumbers(page, totalPages),
	}

	log.Printf("Template data: %+v", data)
	renderTemplate(w, "problemsList.html", data)
}

func generatePageNumbers(currentPage, totalPages int) []int {
	var pages []int
	for i := 1; i <= totalPages; i++ {
		pages = append(pages, i)
	}
	return pages
}

// @desc get HTML page for user's submissions with pagination
// @route GET /submissions?page=<number>
// @access private (only accessible to logged-in users)
func GoSubmissionsPage(w http.ResponseWriter, r *http.Request) {
	userID := "user123" // Replace with actual user ID from auth system

	// Get page number from query parameter, default to 1
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	// Define items per page
	const itemsPerPage = 5

	// TODO: Fetch submissions from database for the current user
	// For now, using static data
	allSubmissions := []Submission{
		{ID: "s1", ProblemID: "1", ProblemTitle: "Two Sum", UserID: userID, Code: "func twoSum(nums []int, target int) []int {\n    for i := 0; i < len(nums); i++ {\n        for j := i + 1; j < len(nums); j++ {\n            if nums[i] + nums[j] == target {\n                return []int{i, j}\n            }\n        }\n    }\n    return nil\n}", Status: "OK", SubmittedAt: time.Now().Add(-24 * time.Hour)},
		{ID: "s2", ProblemID: "2", ProblemTitle: "Add Two Numbers", UserID: userID, Code: "func addTwoNumbers(l1 *ListNode, l2 *ListNode) *ListNode {\n    // ...", Status: "Wrong Answer", SubmittedAt: time.Now().Add(-20 * time.Hour)},
		{ID: "s3", ProblemID: "3", ProblemTitle: "Longest Substring", UserID: userID, Code: "func lengthOfLongestSubstring(s string) int {\n    // ...", Status: "Time Limit", SubmittedAt: time.Now().Add(-15 * time.Hour)},
		{ID: "s4", ProblemID: "4", ProblemTitle: "Three Sum", UserID: userID, Code: "func threeSum(nums []int) [][]int {\n    // ...", Status: "Pending", SubmittedAt: time.Now().Add(-10 * time.Hour)},
		{ID: "s5", ProblemID: "5", ProblemTitle: "Gorg Ali", UserID: userID, Code: "func gorgAli() {\n    // ...", Status: "Compile Error", SubmittedAt: time.Now().Add(-5 * time.Hour)},
		{ID: "s6", ProblemID: "7", ProblemTitle: "DFS", UserID: userID, Code: "func dfs(graph [][]int) {\n    // ...", Status: "OK", SubmittedAt: time.Now()},
	}

	// Filter submissions for the current user (in real app, this would be a DB query)
	var userSubmissions []Submission
	for _, sub := range allSubmissions {
		if sub.UserID == userID {
			userSubmissions = append(userSubmissions, sub)
		}
	}

	// Calculate pagination details
	totalItems := len(userSubmissions)
	totalPages := int(math.Ceil(float64(totalItems) / float64(itemsPerPage)))

	// Ensure page doesn't exceed total pages
	if page > totalPages {
		page = totalPages
	}

	// Calculate start and end indices
	start := (page - 1) * itemsPerPage
	end := start + itemsPerPage
	if end > totalItems {
		end = totalItems
	}

	submissions := userSubmissions[start:end]

	// Generate page numbers for navigation
	var pageNumbers []int
	for i := 1; i <= totalPages; i++ {
		pageNumbers = append(pageNumbers, i)
	}

	// Prepare data for template
	data := struct {
		Submissions []Submission
		CurrentPage int
		PrevPage    int
		NextPage    int
		TotalPages  int
		PageNumbers []int
	}{
		Submissions: submissions,
		CurrentPage: page,
		PrevPage:    page - 1,
		NextPage:    page + 1,
		TotalPages:  totalPages,
		PageNumbers: pageNumbers,
	}

	renderTemplate(w, "my_submission.html", data)
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