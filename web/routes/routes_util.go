package routes

import(
	"net/http"
	"html/template"
)

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	t, err := template.ParseFiles("ui/html/" + tmpl) 
	if err != nil {                                    // Check for errors during template parsing
		http.Error(w, err.Error(), http.StatusInternalServerError) 
	}
	err = t.Execute(w, data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}