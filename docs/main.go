package main

import (
    "fmt"
    "net/http"
)

func main() {
    // Serve static files (CSS and images)
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("templates"))))

    // Serve favicon.ico specifically
    http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "templates/img/favicon.ico")
    })

    // Initialize all page handlers
    InitPage1()
    InitPage2()
    InitPage3()

    // Redirect root to page 1
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        http.Redirect(w, r, "/mypage_example_page1", http.StatusSeeOther)
    })

    // Start the server
    fmt.Println("========================================")
    fmt.Println("HTTPS Server starting on https://mypage.local:8443")
    fmt.Println("========================================")
    fmt.Println("\nIMPORTANT: To use this server, you need to:")
    fmt.Println("1. Add the following line to your hosts file:")
    fmt.Println("   127.0.0.1 mypage.local")
    fmt.Println("\n2. Generate SSL certificates using mkcert:")
    fmt.Println("   mkcert -install")
    fmt.Println("   mkcert mypage.local")
    fmt.Println("\n3. Then access the site at:")
    fmt.Println("   https://mypage.local:8443/mypage_example_page1")
    fmt.Println("\n4. Enter credentials:")
    fmt.Println("   SURNAME: TOCA")
    fmt.Println("   FIRSTNAME: PATRICK")
    fmt.Println("\n========================================")
    
    // Start HTTPS server with mkcert certificates
    if err := http.ListenAndServeTLS(":8443", "mypage.local.pem", "mypage.local-key.pem", nil); err != nil {
        fmt.Printf("Server failed to start: %v\n", err)
    }
}