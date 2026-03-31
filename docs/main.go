package main

import (
    "fmt"
    "html/template"
    "net/http"
)

// Templates
var indexTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>Login Page</title>
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
        }
        button:hover {
            background-color: #45a049;
        }
        .error {
            color: red;
            margin-bottom: 15px;
        }
    </style>
</head>
<body>
    <h2>Enter Your Information</h2>
    {{if .Error}}
    <div class="error">{{.Error}}</div>
    {{end}}
    <form method="POST" action="/">
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
</body>
</html>
`

var secondPageTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>Second Page</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            text-align: center;
            margin-top: 100px;
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
    </style>
</head>
<body>
    <h2>Welcome!</h2>
    <form method="POST" action="/third">
        <button type="submit">Go to Third Page</button>
    </form>
</body>
</html>
`

var thirdPageTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>Third Page</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            text-align: center;
            margin-top: 100px;
        }
        .message {
            font-size: 24px;
            color: #4CAF50;
            margin-bottom: 30px;
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
    </style>
</head>
<body>
    <div class="message">
        <h2>You've reached the third page!</h2>
        <p>Congratulations! 🎉</p>
    </div>
    <form method="GET" action="/">
        <button type="submit">Go Back to Start</button>
    </form>
</body>
</html>
`

func main() {
    // Parse templates
    indexTmpl := template.Must(template.New("index").Parse(indexTemplate))
    secondTmpl := template.Must(template.New("second").Parse(secondPageTemplate))
    thirdTmpl := template.Must(template.New("third").Parse(thirdPageTemplate))

    // Handler for the main page
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == "POST" {
            surname := r.FormValue("surname")
            firstname := r.FormValue("firstname")

            // Check if credentials match TOCA PATRICK
            if surname == "TOCA" && firstname == "PATRICK" {
                // Store in session (using cookie for simplicity)
                http.SetCookie(w, &http.Cookie{
                    Name:  "authenticated",
                    Value: "true",
                    Path:  "/",
                })
                http.Redirect(w, r, "/second", http.StatusSeeOther)
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

    // Handler for second page
    http.HandleFunc("/second", func(w http.ResponseWriter, r *http.Request) {
        // Check if user is authenticated
        cookie, err := r.Cookie("authenticated")
        if err != nil || cookie.Value != "true" {
            http.Redirect(w, r, "/", http.StatusSeeOther)
            return
        }
        secondTmpl.Execute(w, nil)
    })

    // Handler for third page
    http.HandleFunc("/third", func(w http.ResponseWriter, r *http.Request) {
        // Check if user is authenticated
        cookie, err := r.Cookie("authenticated")
        if err != nil || cookie.Value != "true" {
            http.Redirect(w, r, "/", http.StatusSeeOther)
            return
        }

        if r.Method == "POST" {
            thirdTmpl.Execute(w, nil)
            return
        }

        // If not POST, redirect to second page
        http.Redirect(w, r, "/second", http.StatusSeeOther)
    })

    // Start the server
    fmt.Println("Server starting on http://localhost:8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        fmt.Printf("Server failed to start: %v\n", err)
    }
}