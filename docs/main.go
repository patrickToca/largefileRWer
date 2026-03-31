package main

import (
    "fmt"
    "log"
    "net/http"
)

func main() {
    // Serve static files (CSS)
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("templates"))))

    // Initialize all page handlers
    InitPage1()
    InitPage2()
    InitPage3()

    // Redirect root to page 1
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        http.Redirect(w, r, "/mypage_example_page1", http.StatusSeeOther)
    })

    // Optional: Redirect HTTP to HTTPS
    go func() {
        fmt.Println("HTTP server running on :8080 (redirecting to HTTPS)")
        http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            http.Redirect(w, r, "https://mypage.local:8443"+r.RequestURI, http.StatusMovedPermanently)
        }))
    }()

    fmt.Println("========================================")
    fmt.Println("HTTPS Server starting on https://mypage.local:8443")
    fmt.Println("========================================")
    fmt.Println("\nWith mkcert, Chrome will show a SECURE connection!")
    fmt.Println("The padlock icon will appear in the address bar 🔒")
    fmt.Println("\n1. Add the following line to your hosts file:")
    fmt.Println("   127.0.0.1 mypage.local")
    fmt.Println("\n2. Access the site at:")
    fmt.Println("   https://mypage.local:8443/mypage_example_page1")
    fmt.Println("\n3. Enter credentials:")
    fmt.Println("   SURNAME: TOCA")
    fmt.Println("   FIRSTNAME: PATRICK")
    fmt.Println("\n========================================")
    
    // Start HTTPS server with mkcert certificates
    if err := http.ListenAndServeTLS(":8443", "mypage.local.pem", "mypage.local-key.pem", nil); err != nil {
        log.Fatal("HTTPS server failed: ", err)
    }
}