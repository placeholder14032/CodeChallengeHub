package routes
import(
	"net/http"
)

// @desc get landing(welcome) page html
// @route GET /
// @access public
func GoLandingPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "landing.html")
}

// @desc get Html page for auth admin
// @route GET /auth-admin
// @access public
func  GoAuthAdmin(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "adminLogin.html")
}

// @desc get Html page for auth for users
// @route GET /auth-user
// @access public
func  GoAuthUser(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "userLogin.html")
}