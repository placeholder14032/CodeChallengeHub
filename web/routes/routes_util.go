package routes

import(
	"net/http"
	"html/template"
)

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
    t, err := template.ParseFiles("ui/html/" + tmpl)
    if err != nil {
        http.Error(w, "Template not found", http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    if err := t.Execute(w, data); err != nil {
        http.Error(w, "Template rendering failed", http.StatusInternalServerError)
        return
    }
}