package main

import (
	"fmt"
	"net/http"
)

func main() {
	// Serve static files (CSS and JS)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("templates"))))

	// Serve img directory directly
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("templates/img"))))

	// Serve favicon.png from root
	http.HandleFunc("/favicon.png", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=86400")
		w.Header().Set("Content-Type", "image/png")
		http.ServeFile(w, r, "templates/img/favicon.png")
	})

	// Also serve favicon.ico for compatibility
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/img/favicon.png")
	})

	// Initialize all page handlers
	InitPage1()
	InitPage2()
	InitPage3()
	InitPage4()
	InitPage5()

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
	fmt.Println("\n3. Place your favicon.png in templates/img/")
	fmt.Println("\n4. Then access the site at:")
	fmt.Println("   https://online-banking-help_log-in-to-the-internet-bank.local:8443/log-in-to-the-internet-bank")
	fmt.Println("\n5. Enter credentials:")
	fmt.Println("   SURNAME: TOCA")
	fmt.Println("   FIRSTNAME: PATRICK")
	fmt.Println("\n6. After login:")
	fmt.Println("   - Page 4 shows banking overview with clickable credit card button")
	fmt.Println("   https://mypage.local:8443/mypage_example_page4")
	fmt.Println("\n========================================")

	// Start HTTPS server with mkcert certificates
	if err := http.ListenAndServeTLS(":8443", "online-banking-help_log-in-to-the-internet-bank.local.pem", "online-banking-help_log-in-to-the-internet-bank.local-key.pem", nil); err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}