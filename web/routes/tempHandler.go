package routes

import(
	"net/http"
)

func  TempHandler(w http.ResponseWriter, r *http.Request) {
	// Handle the request here
	// w.Write([]byte("Temporary Handler"))
	renderTemplate(w, "temp.html")
}