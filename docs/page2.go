package main

import (
    "html/template"
    "net/http"
    "path/filepath"
)

func InitPage2() {
    // Parse the external HTML template
    tmpl := template.Must(template.ParseFiles(filepath.Join("templates", "page2.html")))
    
    http.HandleFunc("/mypage_example_page2", func(w http.ResponseWriter, r *http.Request) {
        // Check if user is authenticated
        cookie, err := r.Cookie("authenticated")
        if err != nil || cookie.Value != "true" {
            http.Redirect(w, r, "/mypage_example_page1", http.StatusSeeOther)
            return
        }

        // Get username from cookie
        usernameCookie, _ := r.Cookie("username")
        username := "Guest"
        if usernameCookie != nil {
            username = usernameCookie.Value
        }

        // Pass data to template
        data := struct {
            Name string
        }{
            Name: username,
        }
        
        tmpl.Execute(w, data)
    })
}