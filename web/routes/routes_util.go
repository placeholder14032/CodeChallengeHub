package routes

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"log"
)

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
    fmt.Println("ui/html/" + tmpl)
    templatePath := filepath.Join("ui", "html", tmpl)
    t, err := template.ParseFiles(templatePath)
    if err != nil {
        http.Error(w, "Template not found", http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    if err := t.Execute(w, data); err != nil {
        log.Printf("Template rendering failed for %s: %v", tmpl, err) // Add this log
        http.Error(w, "Template rendering failed", http.StatusInternalServerError)
        return
    }
}