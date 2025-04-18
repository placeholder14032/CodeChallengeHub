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



