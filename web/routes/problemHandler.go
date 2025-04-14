package routes

import (
	"fmt"
	"net/http"
	"log"
	"strconv"
	"time"
	// "filepath"	
	"os"

	"github.com/placeHolder143032/CodeChallengeHub/models"
)

// @desc add problem to the database
// @route GET,POST /add_problem
// @access private
func AddProblem(w http.ResponseWriter, r *http.Request) {
	// Parse the HTML template
	tmpl, err := template.ParseFiles("templates/add_problem.html")
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
	userID, ok := r.Context().Value(auth.UserIDKey).(int)
	if !ok {
		log.Println("User ID not found in context")
		tmpl.Execute(w, struct{ Error string; Form FormData }{
			Error: "Authentication error",
			Form:  form,
		})
		return
	}

	// Create problem struct
	problem := models.Problem{
		UserId:        userID,
		Title:         form.Title,
		Difficulty:    form.Difficulty,
		Statement:     form.Statement,
		Input:         form.Input,
		Output:        form.Output,
		CreatedAt:     time.Now(),
		IsPublished:   false, // Draft by default
		TimeLimitMs:   timeLimit,
		MemoryLimitMb: memoryLimit,
	}

	// Save problem to database
	if err := database.AddProblem(userID, problem); err != nil {
		log.Printf("Failed to save problem to database: %v", err)
		tmpl.Execute(w, struct{ Error string; Form FormData }{
			Error: "Failed to save problem",
			Form:  form,
		})
		return
	}

	// Log success
	log.Printf("Problem '%s' created by user %d", form.Title, userID)

	// Redirect to problems list
	http.Redirect(w, r, "/problems", http.StatusSeeOther)
}