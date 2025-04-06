package routes

import(
	"net/http"
	"html/template"
)

func renderTemplate(w http.ResponseWriter, tmpl string) {
	t, err := template.ParseFiles("ui/html/" + tmpl) 
	if err != nil {                                    // Check for errors during template parsing
		http.Error(w, err.Error(), http.StatusInternalServerError) 
	}
	t.Execute(w, nil) 
}