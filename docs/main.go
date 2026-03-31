package main

import (
    "fmt"
    "html/template"
    "net/http"
    "strings"
)

// Templates
var indexTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>MyPage - Example Page 1</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            max-width: 400px;
            margin: 50px auto;
            padding: 20px;
        }
        .form-group {
            margin-bottom: 15px;
        }
        label {
            display: block;
            margin-bottom: 5px;
            font-weight: bold;
        }
        input {
            width: 100%;
            padding: 8px;
            border: 1px solid #ddd;
            border-radius: 4px;
            box-sizing: border-box;
        }
        button {
            background-color: #4CAF50;
            color: white;
            padding: 10px 20px;
            border: none;
            border-radius: 4px;
            cursor: pointer;
            font-size: 16px;
        }
        button:hover {
            background-color: #45a049;
        }
        .error {
            color: red;
            margin-bottom: 15px;
            padding: 10px;
            background-color: #ffebee;
            border-radius: 4px;
        }
        .page-indicator {
            text-align: center;
            color: #666;
            margin-top: 20px;
            font-size: 12px;
        }
    </style>
</head>
<body>
    <h2>Enter Your Information</h2>
    {{if .Error}}
    <div class="error">{{.Error}}</div>
    {{end}}
    <form method="POST" action="/mypage_example_page1">
        <div class="form-group">
            <label for="surname">SURNAME:</label>
            <input type="text" id="surname" name="surname" required>
        </div>
        <div class="form-group">
            <label for="firstname">FIRSTNAME:</label>
            <input type="text" id="firstname" name="firstname" required>
        </div>
        <button type="submit">Submit</button>
    </form>
    <div class="page-indicator">You are on: mypage_example_page1</div>
</body>
</html>
`

var secondPageTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>MyPage - Example Page 2</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            text-align: center;
            margin-top: 100px;
            padding: 20px;
        }
        h2 {
            color: #333;
            margin-bottom: 30px;
        }
        button {
            background-color: #4CAF50;
            color: white;
            padding: 15px 30px;
            font-size: 16px;
            border: none;
            border-radius: 4px;
            cursor: pointer;
        }
        button:hover {
            background-color: #45a049;
        }
        .page-indicator {
            text-align: center;
            color: #666;
            margin-top: 50px;
            font-size: 12px;
        }
    </style>
</head>
<body>
    <h2>Welcome!</h2>
    <form method="POST" action="/mypage_example_page3">
        <button type="submit">Go to Third Page</button>
    </form>
    <div class="page-indicator">You are on: mypage_example_page2</div>
</body>
</html>
`

var thirdPageTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>MyPage - Example Page 3</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            text-align: center;
            margin-top: 100px;
            padding: 20px;
        }
        .message {
            margin-bottom: 30px;
        }
        .message h2 {
            color: #4CAF50;
            font-size: 28px;
            margin-bottom: 10px;
        }
        .message p {
            color: #666;
            font-size: 18px;
        }
        button {
            background-color: #008CBA;
            color: white;
            padding: 10px 20px;
            font-size: 14px;
            border: none;
            border-radius: 4px;
            cursor: pointer;
        }
        button:hover {
            background-color: #007B9E;
        }
        .page-indicator {
            text-align: center;
            color: #666;
            margin-top: 50px;
            font-size: 12px;
        }
    </style>
</head>
<body>
    <div class="message">
        <h2>You've reached the third page!</h2>
        <p>Congratulations! 🎉</p>
    </div>
    <form method="GET" action="/mypage_example_page1">
        <button type="submit">Go Back to Start</button>
    </form>
    <div class="page-indicator">You are on: mypage_example_page3</div>
</body>
</html>
`

func main() {
    // Parse templates
    indexTmpl := template.Must(template.New("index").Parse(indexTemplate))
    secondTmpl := template.Must(template.New("second").Parse(secondPageTemplate))
    thirdTmpl := template.Must(template.New("third").Parse(thirdPageTemplate))

    // Handler for page 1 (login page)
    http.HandleFunc("/mypage_example_page1", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == "POST" {
            surname := r.FormValue("surname")
            firstname := r.FormValue("firstname")

            // Check if credentials match TOCA PATRICK (case-insensitive)
            if strings.ToUpper(surname) == "TOCA" && strings.ToUpper(firstname) == "PATRICK" {
                // Store in session (using cookie for simplicity)
                http.SetCookie(w, &http.Cookie{
                    Name:   "authenticated",
                    Value:  "true",
                    Path:   "/",
                    MaxAge: 3600, // 1 hour
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
            indexTmpl.Execute(w, data)
            return
        }

        // GET request - show login form
        indexTmpl.Execute(w, nil)
    })

    // Handler for page 2
    http.HandleFunc("/mypage_example_page2", func(w http.ResponseWriter, r *http.Request) {
        // Check if user is authenticated
        cookie, err := r.Cookie("authenticated")
        if err != nil || cookie.Value != "true" {
            http.Redirect(w, r, "/mypage_example_page1", http.StatusSeeOther)
            return
        }
        secondTmpl.Execute(w, nil)
    })

    // Handler for page 3
    http.HandleFunc("/mypage_example_page3", func(w http.ResponseWriter, r *http.Request) {
        // Check if user is authenticated
        cookie, err := r.Cookie("authenticated")
        if err != nil || cookie.Value != "true" {
            http.Redirect(w, r, "/mypage_example_page1", http.StatusSeeOther)
            return
        }

        if r.Method == "POST" {
            thirdTmpl.Execute(w, nil)
            return
        }

        // If not POST, redirect to page 2
        http.Redirect(w, r, "/mypage_example_page2", http.StatusSeeOther)
    })

    // Redirect root to the first page
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