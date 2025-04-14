package routes

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"fmt"

	"github.com/placeHolder143032/CodeChallengeHub/database"
	"github.com/placeHolder143032/CodeChallengeHub/middleware"
	"github.com/placeHolder143032/CodeChallengeHub/models"
)

// FormData holds form values for repopulation on error
type FormData struct {
	Title       string
	Difficulty  string
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
	// Parse the HTML template
	tmpl, err := template.ParseFiles("ui/html/add_problem.html")
	if err != nil {
		log.Printf("Template parse error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Handle GET request (render empty form)
	if r.Method == http.MethodGet {
		err := tmpl.Execute(w, nil)
		if err != nil {
			log.Printf("Template execute error: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Handle POST request
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	err = r.ParseForm()
	if err != nil {
		tmpl.Execute(w, struct{ Error string }{Error: "Unable to parse form"})
		return
	}

	// Extract form values
	form := FormData{
		Title:       r.FormValue("title"),
		Difficulty:  r.FormValue("difficulty"),
		Statement:   r.FormValue("statement"),
		TimeLimit:   r.FormValue("time_limit"),
		MemoryLimit: r.FormValue("memory_limit"),
		Input:       r.FormValue("input"),
		Output:      r.FormValue("output"),
	}

	// Validate required fields
	if form.Title == "" || form.Difficulty == "" || form.Statement == "" ||
		form.TimeLimit == "" || form.MemoryLimit == "" || form.Input == "" || form.Output == "" {
		tmpl.Execute(w, struct{ Error string; Form FormData }{
			Error: "All fields are required",
			Form:  form,
		})
		return
	}

	// Validate difficulty
	if form.Difficulty != "Easy" && form.Difficulty != "Medium" && form.Difficulty != "Hard" {
		tmpl.Execute(w, struct{ Error string; Form FormData }{
			Error: "Invalid difficulty level",
			Form:  form,
		})
		return
	}

	// Validate time limit
	timeLimit, err := strconv.Atoi(form.TimeLimit)
	if err != nil || timeLimit <= 0 {
		tmpl.Execute(w, struct{ Error string; Form FormData }{
			Error: "Time limit must be a positive number",
			Form:  form,
		})
		return
	}

	// Validate memory limit
	memoryLimit, err := strconv.Atoi(form.MemoryLimit)
	if err != nil || memoryLimit <= 0 {
		tmpl.Execute(w, struct{ Error string; Form FormData }{
			Error: "Memory limit must be a positive number",
			Form:  form,
		})
		return
	}

	// Get user ID from context
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		log.Println("User ID not found in context")
		tmpl.Execute(w, struct{ Error string; Form FormData }{
			Error: "Authentication error",
			Form:  form,
		})
		return
	}

	// Create directory for problem files
	basePath := "problems"
	timestamp := time.Now().Unix()
	problemDir := filepath.Join(basePath, fmt.Sprintf("problem_%d_%d", userID, timestamp))
	if err := os.MkdirAll(problemDir, 0755); err != nil {
		log.Printf("Failed to create directory %s: %v", problemDir, err)
		tmpl.Execute(w, struct{ Error string; Form FormData }{
			Error: "Failed to save problem files",
			Form:  form,
		})
		return
	}

	// Define file paths
	descriptionPath := filepath.Join(problemDir, "description.txt")
	inputPath := filepath.Join(problemDir, "input.txt")
	outputPath := filepath.Join(problemDir, "output.txt")

	// Save problem content to files
	if err := os.WriteFile(descriptionPath, []byte(form.Statement), 0644); err != nil {
		log.Printf("Failed to write description file %s: %v", descriptionPath, err)
		tmpl.Execute(w, struct{ Error string; Form FormData }{
			Error: "Failed to save problem files",
			Form:  form,
		})
		return
	}

	if err := os.WriteFile(inputPath, []byte(form.Input), 0644); err != nil {
		log.Printf("Failed to write input file %s: %v", inputPath, err)
		tmpl.Execute(w, struct{ Error string; Form FormData }{
			Error: "Failed to save problem files",
			Form:  form,
		})
		return
	}

	if err := os.WriteFile(outputPath, []byte(form.Output), 0644); err != nil {
		log.Printf("Failed to write output file %s: %v", outputPath, err)
		tmpl.Execute(w, struct{ Error string; Form FormData }{
			Error: "Failed to save problem files",
			Form:  form,
		})
		return
	}

	// Create problem struct
	problem := models.Problem{
		UserID:          userID,
		Title:           form.Title,
		// Difficulty:      form.Difficulty, // Uncomment if database supports it
		DescriptionPath: descriptionPath,
		InputPath:       inputPath,
		OutputPath:      outputPath,
		CreatedTime:     time.Now(),
		IsPublished:     false,
		TimeLimit:       timeLimit,
		MemoryLimit:     memoryLimit,
	}

	// Save problem to database
	if err := database.AddProblem(userID, problem); err != nil {
		log.Printf("Failed to save problem to database: %v", err)
		tmpl.Execute(w, struct{ Error string; Form FormData }{
			Error: "Failed to save problem",
			Form:  form,
		})
		// Cleanup files on database failure
		os.RemoveAll(problemDir)
		return
	}

	// Log success
	log.Printf("Problem '%s' created by user %d", form.Title, userID)

	// Redirect to problems list
	http.Redirect(w, r, "/problems", http.StatusSeeOther)
}