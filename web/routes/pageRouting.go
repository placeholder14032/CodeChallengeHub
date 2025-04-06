package routes
import(
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