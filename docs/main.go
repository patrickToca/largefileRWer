package main

import (
    "fmt"
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

    // Start the server
    fmt.Println("========================================")
    fmt.Println("Server starting on http://mypage.local:8080")
    fmt.Println("========================================")
    fmt.Println("\nIMPORTANT: To use this server, you need to:")
    fmt.Println("1. Add the following line to your hosts file:")
    fmt.Println("   127.0.0.1 mypage.local")
    fmt.Println("\n   On Linux/Mac: /etc/hosts")
    fmt.Println("   On Windows: C:\\Windows\\System32\\drivers\\etc\\hosts")
    fmt.Println("\n2. Then access the site at:")
    fmt.Println("   http://mypage.local:8080/mypage_example_page1")
    fmt.Println("\n3. Enter credentials:")
    fmt.Println("   SURNAME: TOCA")
    fmt.Println("   FIRSTNAME: PATRICK")
    fmt.Println("\n========================================")
    
    if err := http.ListenAndServe(":8080", nil); err != nil {
        fmt.Printf("Server failed to start: %v\n", err)
    }
}