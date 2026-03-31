package main

import (
    "html/template"
    "net/http"
    "path/filepath"
)

func InitPage3() {
    // Parse the external HTML template
    tmpl := template.Must(template.ParseFiles(filepath.Join("templates", "page3.html")))
    
    http.HandleFunc("/mypage_example_page3", func(w http.ResponseWriter, r *http.Request) {
        // Check if user is authenticated
        cookie, err := r.Cookie("authenticated")
        if err != nil || cookie.Value != "true" {
            http.Redirect(w, r, "/mypage_example_page1", http.StatusSeeOther)
            return
        }

        if r.Method == "POST" {
            tmpl.Execute(w, nil)
            return
        }

        // If not POST, redirect to page 2
        http.Redirect(w, r, "/mypage_example_page2", http.StatusSeeOther)
    })
}