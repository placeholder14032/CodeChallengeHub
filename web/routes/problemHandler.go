package routes

import (
	"net/http"
	"fmt"
)

// @desc submit answer for problem
// @route POST /api/submit_answer 
// @access private (you can only access this page if you are logged in) 
func SubmitAnswer(w http.ResponseWriter, r *http.Request) {
	fmt.Print("Submit Answer")
 }