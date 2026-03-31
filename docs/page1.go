package main

import (
    "html/template"
    "net/http"
    "path/filepath"
    "strings"
)

func InitPage1() {
    // Parse the external HTML template
    tmpl := template.Must(template.ParseFiles(filepath.Join("templates", "html", "page1.html")))
    
    http.HandleFunc("/mypage_example_page1", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == "POST" {
            surname := r.FormValue("surname")
            firstname := r.FormValue("firstname")

            // Check if credentials match TOCA PATRICK (case-insensitive)
            if strings.ToUpper(surname) == "TOCA" && strings.ToUpper(firstname) == "PATRICK" {
                // Store user info in session (using cookie for simplicity)
                http.SetCookie(w, &http.Cookie{
                    Name:   "authenticated",
                    Value:  "true",
                    Path:   "/",
                    MaxAge: 3600, // 1 hour
                })
                http.SetCookie(w, &http.Cookie{
                    Name:   "username",
                    Value:  firstname,
                    Path:   "/",
                    MaxAge: 3600,
                })
                http.Redirect(w, r, "/mypage_example_page2", http.StatusSeeOther)
                return
            }

            // If authentication fails, show error
            data := struct {
                Error string
            }{
                Error: "Invalid credentials. Please try again.",
            }
            tmpl.Execute(w, data)
            return
        }

        // GET request - show login form
        tmpl.Execute(w, nil)
    })
}