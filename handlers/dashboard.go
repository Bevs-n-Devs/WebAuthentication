package handlers

import (
	"fmt"
	"net/http"

	"github.com/Bevs-n-Devs/WebAuthentication/logs"
	"github.com/Bevs-n-Devs/WebAuthentication/middleware"
)

func Dashboard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		logs.Logs(logWarning, fmt.Sprintf("Invalid request method: %s. Redirecting back to login page...", r.Method))
		http.Redirect(w, r, "/login", http.StatusBadRequest)
		return
	}

	// denies the request if authorization fails
	err := middleware.AuthorizeRequest(r)
	if err != nil {
		logs.Logs(logWarning, fmt.Sprintf("Failed to authorize request: %s. Redirecting back to login page...", err.Error()))
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
		return
	}

	// direct user to protected page after authorization
	err = Templates.ExecuteTemplate(w, "dashboard.html", nil)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Failed to execute template: %s", err.Error()))
		http.Error(w, fmt.Sprintf("Unable to load page: %s", err.Error()), http.StatusInternalServerError)
	}

}

// func Dashboard(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodGet {
// 		logs.Logs(logWarning, fmt.Sprintf("Invalid request method: %s. Redirecting back to login page...", r.Method))
// 		http.Redirect(w, r, "/login", http.StatusBadRequest)
// 		return
// 	}

// 	// denies the request if authorization fails
// 	err := middleware.AuthorizeRequest(r)
// 	if err != nil {
// 		logs.Logs(logWarning, fmt.Sprintf("Failed to authorize request: %s. Redirecting back to login page...", err.Error()))
// 		http.Redirect(w, r, "/login", http.StatusUnauthorized)
// 		return
// 	}

// 	// get the username from the session cookie
// 	sessionCookie, err := r.Cookie("session_token")
// 	if err != nil {
// 		logs.Logs(logWarning, "Session cookie is missing")
// 		http.Redirect(w, r, "/login", http.StatusUnauthorized)
// 		return
// 	}

// 	// verify the session token
// 	ok, err := db.ValidateSessionToken(sessionCookie.Value)
// 	if err != nil {
// 		logs.Logs(logErr, fmt.Sprintf("Failed to validate session token: %s", err.Error()))
// 		http.Error(w, fmt.Sprintf("Unable to load page: %s", err.Error()), http.StatusInternalServerError)
// 		return
// 	}

// 	if !ok {
// 		logs.Logs(logWarning, "Invalid session token")
// 		http.Redirect(w, r, "/login", http.StatusUnauthorized)
// 		return
// 	}

// 	// get the username from the database
// 	username, err := db.GetUsernameFromSessionToken(sessionCookie.Value)
// 	if err != nil {
// 		logs.Logs(logErr, fmt.Sprintf("Failed to get username from session token: %s", err.Error()))
// 		http.Error(w, fmt.Sprintf("Unable to load page: %s", err.Error()), http.StatusInternalServerError)
// 		return
// 	}

// 	logs.Logs(logInfo, fmt.Sprintf("User %s has successfully logged in. Directed to dashboard page...", username))

// 	err = Templates.ExecuteTemplate(w, "dashboard.html", nil)
// 	if err != nil {
// 		logs.Logs(logErr, fmt.Sprintf("Failed to execute template: %s", err.Error()))
// 		http.Error(w, fmt.Sprintf("Unable to load page: %s", err.Error()), http.StatusInternalServerError)
// 	}
// }
