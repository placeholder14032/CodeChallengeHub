package routes

import (
	"fmt"
	"net/http"
	"log"
	"strconv"
	"time"
)

// @desc submit answer for problem
// @route POST /api/submit_answer
// @access private (you can only access this page if you are logged in)
func SubmitAnswer(w http.ResponseWriter, r *http.Request) {
	fmt.Print("Submit Answer")
 }

// struct for now later we should make it compatuble with other team "khoshnevis i guess"
type Problem struct {
    ID           string
    Title        string
    Difficulty   string
    Solved       bool
    Submissions  int
    Statement    string
    TimeLimit    int
    MemoryLimit  int
    Input        string
    Output       string
    OwnerID      string
    Published    bool
}
// @desc add problem to the database
// @route POST /api/add_problem
// @access private (any logged-in user can create a draft)
func AddProblem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Assume user is authenticated
	userID := "user123" // TODO: Replace with actual user ID from session/token

	// Parse form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Extract form values
	title := r.FormValue("title")
	difficulty := r.FormValue("difficulty")
	statement := r.FormValue("statement")
	timeLimitStr := r.FormValue("time_limit")
	memoryLimitStr := r.FormValue("memory_limit")
	input := r.FormValue("input")
	output := r.FormValue("output")

	// Validate inputs
	if title == "" || difficulty == "" || statement == "" || timeLimitStr == "" || memoryLimitStr == "" || input == "" || output == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	timeLimit, err := strconv.Atoi(timeLimitStr)
	if err != nil || timeLimit <= 0 {
		http.Error(w, "Invalid time limit", http.StatusBadRequest)
		return
	}

	memoryLimit, err := strconv.Atoi(memoryLimitStr)
	if err != nil || memoryLimit <= 0 {
		http.Error(w, "Invalid memory limit", http.StatusBadRequest)
		return
	}

	// Create a new problem (draft)
	problem := Problem{
		ID:          "p" + strconv.Itoa(int(time.Now().Unix())), // Simple ID generation
		Title:       title,
		Difficulty:  difficulty,
		Solved:      false,
		Submissions: 0,
		Statement:   statement,
		TimeLimit:   timeLimit,
		MemoryLimit: memoryLimit,
		Input:       input,
		Output:      output,
		OwnerID:     userID,
		Published:   false, // Starts as draft
	}

	// TODO: Save problem to database
	log.Printf("New problem created: %+v", problem)

	// Redirect to problems list or a "My Problems" page
	http.Redirect(w, r, "/problems", http.StatusSeeOther)
}