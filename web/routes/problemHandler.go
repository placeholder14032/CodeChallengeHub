package routes

import (
	"fmt"
	"net/http"
)

// @desc submit answer for problem
// @route POST /api/submit_answer
// @access private (you can only access this page if you are logged in)
func SubmitAnswer(w http.ResponseWriter, r *http.Request) {
	fmt.Print("Submit Answer")
 }

// struct for now later we should make it compatuble with other team "khoshnevis i guess"
 type Problem struct {
    ID          string
    Title       string
    Difficulty  string
    Solved      bool
    Submissions int
}

type PaginatedProblems struct {
    Problems     []Problem
    CurrentPage  int
    TotalPages   int
    HasNext      bool
    HasPrev      bool
    ItemsPerPage int
}

