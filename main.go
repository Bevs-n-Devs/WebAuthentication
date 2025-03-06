package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type Login struct {
	HashedPassword string
	SessionToken   string
	CSRFToken      string
}

// mock database (in memory)
var usersDB = map[string]Login{}

func main() {
	fmt.Println("Hello world, hello Yaw!")
	log.Println("Starting Web Authentication web app...")
	http.HandleFunc("/register", register)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/protected", protected)
	log.Println("Starting web app on port 8080")
	http.ListenAndServe(":8080", nil)
}

func register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	// get the username & password from form
	username := r.FormValue("username")
	password := r.FormValue("password")
	// basic validation
	if len(username) == 0 || len(password) == 0 {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// check if the user already exists
	if _, ok := usersDB[username]; ok {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	// hash password & store in mock database
	hashedPassword, err := HashedPassword(password)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	usersDB[username] = Login{
		HashedPassword: hashedPassword,
	}

	fmt.Fprintln(w, "User registered successfully")

}

// here we use session & SRF token to login user
func login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	// check if the user exists
	user, ok := usersDB[username]
	if !ok || !checkPasswordHash(password, user.HashedPassword) {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	sessionToken := generateToken(32)
	csrfToken := generateToken(32)

	// set session cookie to send back to the client
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  time.Now().Add(24 * time.Hour), // expire in 24 hours
		HttpOnly: true,
	})

	// set CSRF token in a cookie to send back to the client
	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    csrfToken,
		Expires:  time.Now().Add(24 * time.Hour), // expire in 24 hours
		HttpOnly: false,                          // needs to be accessible to the client so can retrieve it from header
	})

	// store session token in mock database
	user.SessionToken = sessionToken // store session token
	user.CSRFToken = csrfToken       // store CSRF token
	usersDB[username] = user         // update mock database

	fmt.Fprintf(w, "Login successful! Welcome %s!", username)
}

// implementing teh protected route
func protected(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// this denies the request if error is recieved
	if err := Authorize(r); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	username := r.FormValue("username")
	fmt.Fprintf(w, "CSRF token is valid! Welcome %s!", username)
}

func logout(w http.ResponseWriter, r *http.Request) {
	if err := Authorize(r); err != nil {
		http.Error(w, fmt.Sprintf("Unauthorized: %s", err.Error()), http.StatusUnauthorized)
		return
	}

	// clear cookie
	http.SetCookie(w, &http.Cookie{
		Name:  "session_token",
		Value: "",
		// Immediately expire by setting expiration time to one hour ago
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
	})
	http.SetCookie(w, &http.Cookie{
		Name:  "csrf_token",
		Value: "",
		// Immediately expire by setting expiration time to one hour ago
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: false,
	})

	// clear the tokens from the mock database
	username := r.FormValue("username")
	if user, ok := usersDB[username]; ok {
		user.SessionToken = ""
		user.CSRFToken = ""
		usersDB[username] = user
	}

	fmt.Fprintln(w, "Logout successful!")
}
