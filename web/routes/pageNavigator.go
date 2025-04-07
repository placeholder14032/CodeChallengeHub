package routes
import(
	"net/http"
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