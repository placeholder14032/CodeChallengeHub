package routes

import (
    "net/http"
    "html/template"
    "strings"
    "log"
    "path/filepath"
)

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
    // Ensure the template name includes the .html extension
    if !strings.HasSuffix(tmpl, ".html") {
        tmpl += ".html"
    }

    // Define the custom function for replacing spaces in template
    funcMap := template.FuncMap{
        "replaceSpaces": func(s string) string {
            return strings.ReplaceAll(s, " ", "-")
        },
    }

    // Parse the template file
    templatePath := filepath.Join("ui/html", tmpl)
    t, err := template.New(filepath.Base(templatePath)).Funcs(funcMap).ParseFiles(templatePath)
    if err != nil {
        log.Printf("Failed to parse template %s: %v", templatePath, err)
        http.Error(w, "Failed to load template: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Set the Content-Type header
    w.Header().Set("Content-Type", "text/html; charset=utf-8")

    // Execute the template with the provided data
    if err := t.Execute(w, data); err != nil {
        log.Printf("Template rendering failed for %s: %v", templatePath, err)
        http.Error(w, "Failed to render template: "+err.Error(), http.StatusInternalServerError)
        return
    }
}