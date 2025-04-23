package routes

import (
	"fmt"
	"html/template"
	"log"
	// "math"
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
	// log.Printf("AddProblem: Request received, Method: %s, URL: %s", r.Method, r.URL.Path)

	// Parse the HTML template
	tmpl, err := template.ParseFiles("ui/html/add_problem.html")
	if err != nil {
		// log.Printf("AddProblem: Template parse error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	// log.Printf("AddProblem: Template parsed successfully")

	// Handle GET request (render empty form)
	if r.Method == http.MethodGet {
		// log.Printf("AddProblem: Handling GET request")
		err := tmpl.Execute(w, nil)
		if err != nil {
			// log.Printf("AddProblem: Template execute error: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		} else {
			// log.Printf("AddProblem: GET request served, form rendered")
		}
		return
	}

	// Handle POST request
	if r.Method != http.MethodPost {
		// log.Printf("AddProblem: Invalid method %s, expected POST", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// log.Printf("AddProblem: Handling POST request")

	// Parse form data
	err = r.ParseForm()
	if err != nil {
		// log.Printf("AddProblem: Form parse error: %v", err)
		tmpl.Execute(w, struct{ Error string }{Error: "Unable to parse form"})
		return
	}
	// log.Printf("AddProblem: Form parsed successfully")

	// Extract form values
	form := FormData{
		Title:       r.FormValue("title"),
		Statement:   r.FormValue("statement"),
		TimeLimit:   r.FormValue("time_limit"),
		MemoryLimit: r.FormValue("memory_limit"),
		Input:       r.FormValue("input"),
		Output:      r.FormValue("output"),
	}
	// log.Printf("AddProblem: Form data extracted: Title=%s, TimeLimit=%s, MemoryLimit=%s",
		// form.Title, form.TimeLimit, form.MemoryLimit)

	// Validate required fields
	if form.Title == "" || form.Statement == "" ||
		form.TimeLimit == "" || form.MemoryLimit == "" || form.Input == "" || form.Output == "" {
		// log.Printf("AddProblem: Validation failed: Missing required fields")
		tmpl.Execute(w, struct{ Error string; Form FormData }{
			Error: "All fields are required",
			Form:  form,
		})
		return
	}
	// log.Printf("AddProblem: Required fields validated")

	// Validate time limit
	timeLimit, err := strconv.Atoi(form.TimeLimit)
	if err != nil || timeLimit <= 0 {
		// log.Printf("AddProblem: Validation failed: Invalid time limit: %s, error: %v", form.TimeLimit, err)
		tmpl.Execute(w, struct{ Error string; Form FormData }{
			Error: "Time limit must be a positive number",
			Form:  form,
		})
		return
	}
	// log.Printf("AddProblem: Time limit validated: %d ms", timeLimit)

	// Validate memory limit
	memoryLimit, err := strconv.Atoi(form.MemoryLimit)
	if err != nil || memoryLimit <= 0 {
		// log.Printf("AddProblem: Validation failed: Invalid memory limit: %s, error: %v", form.MemoryLimit, err)
		tmpl.Execute(w, struct{ Error string; Form FormData }{
			Error: "Memory limit must be a positive number",
			Form:  form,
		})
		return
	}
	// log.Printf("AddProblem: Memory limit validated: %d MB", memoryLimit)

	// Get user ID from context
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		// log.Printf("AddProblem: User ID not found in context")
		tmpl.Execute(w, struct{ Error string; Form FormData }{
			Error: "Authentication error",
			Form:  form,
		})
		return
	}
	// log.Printf("AddProblem: User ID retrieved: %d", userID)

	// Create directory for problem files
	basePath := "pkg"
	userDir := fmt.Sprintf("%v", userID)
	problemDir := filepath.Join(basePath, userDir,"problemsCreated", fmt.Sprintf(form.Title))

	// log.Printf("AddProblem: Creating directory: %s", problemDir)
	if err := os.MkdirAll(problemDir, 0755); err != nil {
		// log.Printf("AddProblem: Failed to create directory %s: %v", problemDir, err)
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
	log.Printf("AddProblem: Redirecting to /allproblems-user")
	http.Redirect(w, r, "/allproblems-user", http.StatusSeeOther)
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
// func GoSubmissionsPage(w http.ResponseWriter, r *http.Request) {
// 	userID := "user123" // Replace with actual user ID from auth system

// 	// Get page number from query parameter, default to 1
// 	pageStr := r.URL.Query().Get("page")
// 	page, err := strconv.Atoi(pageStr)
// 	if err != nil || page < 1 {
// 		page = 1
// 	}

// 	// Define items per page
// 	const itemsPerPage = 5

// 	// TODO: Fetch submissions from database for the current user
// 	// For now, using static data
// 	allSubmissions := []Submission{
// 		{ID: "s1", ProblemID: "1", ProblemTitle: "Two Sum", UserID: userID, Code: "func twoSum(nums []int, target int) []int {\n    for i := 0; i < len(nums); i++ {\n        for j := i + 1; j < len(nums); j++ {\n            if nums[i] + nums[j] == target {\n                return []int{i, j}\n            }\n        }\n    }\n    return nil\n}", Status: "OK", SubmittedAt: time.Now().Add(-24 * time.Hour)},
// 		{ID: "s2", ProblemID: "2", ProblemTitle: "Add Two Numbers", UserID: userID, Code: "func addTwoNumbers(l1 *ListNode, l2 *ListNode) *ListNode {\n    // ...", Status: "Wrong Answer", SubmittedAt: time.Now().Add(-20 * time.Hour)},
// 		{ID: "s3", ProblemID: "3", ProblemTitle: "Longest Substring", UserID: userID, Code: "func lengthOfLongestSubstring(s string) int {\n    // ...", Status: "Time Limit", SubmittedAt: time.Now().Add(-15 * time.Hour)},
// 		{ID: "s4", ProblemID: "4", ProblemTitle: "Three Sum", UserID: userID, Code: "func threeSum(nums []int) [][]int {\n    // ...", Status: "Pending", SubmittedAt: time.Now().Add(-10 * time.Hour)},
// 		{ID: "s5", ProblemID: "5", ProblemTitle: "Gorg Ali", UserID: userID, Code: "func gorgAli() {\n    // ...", Status: "Compile Error", SubmittedAt: time.Now().Add(-5 * time.Hour)},
// 		{ID: "s6", ProblemID: "7", ProblemTitle: "DFS", UserID: userID, Code: "func dfs(graph [][]int) {\n    // ...", Status: "OK", SubmittedAt: time.Now()},
// 	}

// 	// Filter submissions for the current user (in real app, this would be a DB query)
// 	var userSubmissions []Submission
// 	for _, sub := range allSubmissions {
// 		if sub.UserID == userID {
// 			userSubmissions = append(userSubmissions, sub)
// 		}
// 	}

// 	// Calculate pagination details
// 	totalItems := len(userSubmissions)
// 	totalPages := int(math.Ceil(float64(totalItems) / float64(itemsPerPage)))

// 	// Ensure page doesn't exceed total pages
// 	if page > totalPages {
// 		page = totalPages
// 	}

// 	// Calculate start and end indices
// 	start := (page - 1) * itemsPerPage
// 	end := start + itemsPerPage
// 	if end > totalItems {
// 		end = totalItems
// 	}

// 	submissions := userSubmissions[start:end]

// 	// Generate page numbers for navigation
// 	var pageNumbers []int
// 	for i := 1; i <= totalPages; i++ {
// 		pageNumbers = append(pageNumbers, i)
// 	}

// 	// Prepare data for template
// 	data := struct {
// 		Submissions []Submission
// 		CurrentPage int
// 		PrevPage    int
// 		NextPage    int
// 		TotalPages  int
// 		PageNumbers []int
// 	}{
// 		Submissions: submissions,
// 		CurrentPage: page,
// 		PrevPage:    page - 1,
// 		NextPage:    page + 1,
// 		TotalPages:  totalPages,
// 		PageNumbers: pageNumbers,
// 	}

// 	renderTemplate(w, "my_submission.html", data)
// }

// @desc get HTML page for each problem page
// @route GET /problem
// @access private (you can only access this page if you are logged in)
func GoProblemPage(w http.ResponseWriter, r *http.Request) {
    // Get problem ID from URL parameter
    problemID := r.URL.Query().Get("id")
    if problemID == "" {
        http.Error(w, "Problem ID is required", http.StatusBadRequest)
        return
    }

    // Convert problemID to int
    id, err := strconv.Atoi(problemID)
    if err != nil {
        http.Error(w, "Invalid problem ID", http.StatusBadRequest)
        return
    }

    // Get userID from context with detailed logging
    userIDValue := r.Context().Value(middleware.UserIDKey)
    if userIDValue == nil {
        log.Printf("Unauthorized: No UserIDKey in context for request %s", r.URL.String())
        http.Redirect(w, r, "/login-user", http.StatusSeeOther)
        return
    }

    userID, ok := userIDValue.(int)
    if !ok {
        log.Printf("Unauthorized: UserIDKey value is not an int, got type %T, value %v, for request %s", 
            userIDValue, userIDValue, r.URL.String())
        http.Redirect(w, r, "/login-user", http.StatusSeeOther)
        return
    }

    // Get problem from database
    problem, err := database.GetSingleProblem(id)
    if err != nil {
        if err.Error() == "problem does not exist" {
            http.Error(w, "Problem not found", http.StatusNotFound)
            return
        }
        log.Printf("Error fetching problem %d: %v", id, err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    // isAdmin, err := database.GetUserRole(userID)
    if err != nil {
        log.Printf("Error fetching user role for user %d: %v", userID, err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    // if !problem.IsPublished && isAdmin == 0 {
    //     http.Error(w, "Problem not available", http.StatusForbidden)
    //     return
    // }

    // Read description, input, and output files
    description, err := os.ReadFile(problem.DescriptionPath)
    if err != nil {
        log.Printf("Error reading description file %s: %v", problem.DescriptionPath, err)
        http.Error(w, "Error reading problem description", http.StatusInternalServerError)
        return
    }

    input, err := os.ReadFile(problem.InputPath)
    if err != nil {
        log.Printf("Error reading input file %s: %v", problem.InputPath, err)
        http.Error(w, "Error reading input description", http.StatusInternalServerError)
        return
    }

    output, err := os.ReadFile(problem.OutputPath)
    if err != nil {
        log.Printf("Error reading output file %s: %v", problem.OutputPath, err)
        http.Error(w, "Error reading output description", http.StatusInternalServerError)
        return
    }

    // Prepare template data
    data := models.ProblemData{
        Title:     problem.Title,
        Explanation: string(description),
        Input:       string(input),
        Output:      string(output),
		TimeLimit:   problem.TimeLimit,
		MemoryLimit: problem.MemoryLimit,
		ID:          problem.ID,
    }

    w.Header().Set("Content-Type", "text/html; charset=utf-8")

    renderTemplate(w, "problem.html", data)
}