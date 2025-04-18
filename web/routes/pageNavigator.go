package routes

import (
    "net/http"
)

// @desc get landing(welcome) page html
// @route GET /
// @access public
func GoLandingPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "landing.html", nil)
}

// @desc get Html page for auth admin
// @route GET /login-admin
// @access public
func GoLoginAdmin(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "auth/adminLogin.html", nil)
}

// @desc get Html page for login for users
// @route GET /login-user
// @access public
func GoLoginUser(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "auth/userLogin.html", nil)
}

// @desc get Html page for auth admin
// @route GET /signup-admin
// @access public
func GoSignupAdmin(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "auth/adminSignup.html", nil)
}

// @desc get Html page for login for users
// @route GET /signup-user
// @access public
func GoSignupUser(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "auth/userSignup.html", nil)
}

// @desc get Html page for each problem page
// @route GET /problem -> // later we should search with id i guess: /problem/:id ??
// @access private (you can only access this page if you are logged in)
func GoProblemPage(w http.ResponseWriter, r *http.Request) {
	problem := struct {
		ID          string
		Title       string
		Statement   string
		Explanation string
		TimeLimit   int
		MemoryLimit int
		Input       string
		Output      string
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

	renderTemplate(w, "problem.html", problem)
}

