package routes

import (
	"math"
	"net/http"
	"strconv"
	"time"
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


// @desc get HTML page for adding a new problem
// @route GET /add_problem
// @access private (only logged-in users)
func GoAddProblemPage(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        renderTemplate(w, "add_problem.html", nil)
        return
    }
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