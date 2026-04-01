package main

import (
    "html/template"
    "net/http"
    "path/filepath"
)

func InitPage4() {
    // Parse the external HTML template
    tmpl := template.Must(template.ParseFiles(filepath.Join("templates", "html", "page4.html")))
    
    http.HandleFunc("/mypage_example_page4", func(w http.ResponseWriter, r *http.Request) {
        // Check if user is authenticated
        cookie, err := r.Cookie("authenticated")
        if err != nil || cookie.Value != "true" {
            http.Redirect(w, r, "/mypage_example_page1", http.StatusSeeOther)
            return
        }

        // Execute the template
        tmpl.Execute(w, nil)
    })
}