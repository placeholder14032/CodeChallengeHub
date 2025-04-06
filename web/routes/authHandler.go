package routes

import (
	"fmt"
	"net/http"
)

// @desc get landing(welcome) page html
// @route GET /
// @access public
func  LandingHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "landing.html")
}

// @desc get Html page for auth
// @route GET /auth
// @access public
func  AuthHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "auth.html")
}

// @desc login and check password for signin user
// @route POST /api/auth/login
// @access public
func  LoginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Print("Loginnnnnnn")
}

// @desc creaing a new account for signup
// @route POST /api/auth/register
// @access public
func  SigninHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Print("Signinnnnnnnn")
}