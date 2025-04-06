package routes

import(
	"net/http"
)

func  LandingHandler(w http.ResponseWriter, r *http.Request) {
	// Handle the request here
	// w.Write([]byte("Temporary Handler"))
	renderTemplate(w, "landing.html")
}

func  AuthHandler(w http.ResponseWriter, r *http.Request) {
	// Handle the request here
	// w.Write([]byte("Temporary Handler"))
	renderTemplate(w, "auth.html")
}