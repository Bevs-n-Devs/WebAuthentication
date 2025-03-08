package handlers

import (
	"fmt"
	"net/http"

	"github.com/Bevs-n-Devs/WebAuthentication/db"
	"github.com/Bevs-n-Devs/WebAuthentication/logs"
)

func SubmitLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logs.Logs(logWarning, fmt.Sprintf("Invalid request method: %s. Redirecting back to index page...", r.Method))
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// parse form data
	err := r.ParseForm()
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Failed to parse form data: %s", err.Error()))
		http.Error(w, fmt.Sprintf("Unable to load page: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// get form data
	username := r.FormValue("username")
	password := r.FormValue("password")

	// check if user exists in database
	exists, err := db.AuthenticateUser(username, password)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Failed to authenticate user: %s", err.Error()))
		http.Error(w, fmt.Sprintf("Unable to load page: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	if !exists {
		logs.Logs(logWarning, "User does not exist or invalid password. Redirecting back to login page...")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// add tokens to user in database
	sessionToken, csrfToken, expiry, err := db.UpdateSessionTokens(username)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Failed to update session tokens: %s", err.Error()))
		http.Error(w, fmt.Sprintf("Unable to load page: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// set session cookie for client
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  expiry, // exires in 24hrs (same as database expiry)
		HttpOnly: true,
	})

	// set CSRF token in a cookie for client
	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    csrfToken,
		Expires:  expiry, // exires in 24hrs (same as database expiry)
		HttpOnly: false,  // allows client to access CSRF token
	})

	// redirect to dashboard page if authentication is successful
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}
