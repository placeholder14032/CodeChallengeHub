package routes

import (
	"math"
	"net/http"
	"strconv"
)

// @desc get landing(welcome) page html
// @route GET /
// @access public
func GoLandingPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "landing.html",nil)
}

// @desc get Html page for auth admin
// @route GET /login-admin
// @access public
func  GoLoginAdmin(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "auth/adminLogin.html",nil)
}

// @desc get Html page for login for users
// @route GET /login-user
// @access public
func  GoLoginUser(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "auth/userLogin.html",nil)
}

// @desc get Html page for auth admin
// @route GET /signup-admin
// @access public
func  GoSignupAdmin(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "auth/adminSignup.html",nil)
}

// @desc get Html page for login for users
// @route GET /signup-user
// @access public
func  GoSignupUser(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "auth/userSignup.html",nil)
}

// @desc get Html page for profile page
// @route GET /profile
// @access private (you can only access this page if you are logged in)
func  GoProfilePage(w http.ResponseWriter, r *http.Request) {

	data := struct {
		Username        string
		TotalAttempts   int
		SuccessRate     float64
		QuestionsSolved int
		IsAdmin         bool
		UserID          string
	}{
		Username:        "Redsunnn",
		TotalAttempts:   91,
		SuccessRate:     75.5,
		QuestionsSolved: 75,
		IsAdmin:         false, 
		UserID:          "123",
	}

	renderTemplate(w, "profilePage.html",data)
}


// @desc get Html page for each problem page
// @route GET /problem -> // later we should search with id i guess: /problem/:id ??
// @access private (you can only access this page if you are logged in) 
func GoProblemPage(w http.ResponseWriter, r *http.Request) {
	 problem := struct{
		ID         string
		Title      string
		Statement  string
		Explanation string
		TimeLimit  int
		MemoryLimit int
		Input     string
		Output     string   
	 }{
        ID:          "1",
        Title:       "Add Two Numbers",
        Statement:   "Write a program to add two numbers and return their sum.",
        Explanation: "You are given two integers as input. Your task is to compute their sum and output the result.",
        TimeLimit:   1000,
        MemoryLimit: 256,
        Input:       "2 3",
        Output:      "5",
    }


	renderTemplate(w, "problem.html",problem)
}

// @desc get Html page for submitting problem
// @route GET / submit_answer 
// @access private (you can only access this page if you are logged in) 
func GoSubmitAnswer(w http.ResponseWriter, r *http.Request) {
   renderTemplate(w, "problem_submit.html",nil)
}

// @desc get HTML page for all problems with pagination
// @route GET /problems?page=<number>
// @access private (you can only access this page if you are logged in)
func GoProblemsListPage(w http.ResponseWriter, r *http.Request) {
    // Get page number from query parameter, default to 1
    pageStr := r.URL.Query().Get("page")
    page, err := strconv.Atoi(pageStr)
    if err != nil || page < 1 {
        page = 1
    }

    // Define items per page
    const itemsPerPage = 5

    // TODO: Fetch problems from your database
    // For now, using static data
    allProblems := []Problem{
        {ID: "1", Title: "Two Sum", Difficulty: "Easy", Solved: false, Submissions: 0},
        {ID: "2", Title: "Add Two Numbers", Difficulty: "Medium", Solved: false, Submissions: 0},
        {ID: "3", Title: "Longest Substring", Difficulty: "Medium", Solved: false, Submissions: 0},
        {ID: "4", Title: "Three Sum", Difficulty: "Easy", Solved: false, Submissions: 0},
        {ID: "5", Title: "Gorg Ali", Difficulty: "Medium", Solved: false, Submissions: 0},
        {ID: "7", Title: "DFS", Difficulty: "Medium", Solved: false, Submissions: 0},
        {ID: "8", Title: "A*", Difficulty: "Easy", Solved: false, Submissions: 0},
        {ID: "9", Title: "Othello", Difficulty: "Medium", Solved: false, Submissions: 0},
        {ID: "10", Title: "Project", Difficulty: "Medium", Solved: false, Submissions: 0},
        {ID: "11", Title: "Super Hexagon", Difficulty: "Easy", Solved: false, Submissions: 0},
        {ID: "12", Title: "Super Mario", Difficulty: "Medium", Solved: false, Submissions: 0},
        {ID: "13", Title: "Justhis", Difficulty: "Medium", Solved: false, Submissions: 0},
        {ID: "14", Title: "Big Mouth", Difficulty: "Easy", Solved: false, Submissions: 0},
        {ID: "15", Title: "New Face", Difficulty: "Medium", Solved: false, Submissions: 0},
        {ID: "16", Title: "Alsfnnlsnfcnasf", Difficulty: "Medium", Solved: false, Submissions: 0},
    }

    // Calculate pagination details
    totalItems := len(allProblems)
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

    // Slice the problems for the current page
    problems := allProblems[start:end]

    // Generate page numbers for navigation (simple version: show all pages)
    var pageNumbers []int
    for i := 1; i <= totalPages; i++ {
        pageNumbers = append(pageNumbers, i)
    }

    // Prepare data for template
    data := struct {
        Problems    []Problem
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
        PageNumbers: pageNumbers,
    }

    renderTemplate(w, "problemsList.html", data)
}