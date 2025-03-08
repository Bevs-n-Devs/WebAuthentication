package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Bevs-n-Devs/WebAuthentication/db"
	"github.com/Bevs-n-Devs/WebAuthentication/logs"
	"github.com/Bevs-n-Devs/WebAuthentication/middleware"
)

func LogoutUser(w http.ResponseWriter, r *http.Request) {
	err := middleware.AuthorizeRequest(r)
	if err != nil {
		logs.Logs(logWarning, fmt.Sprintf("Failed to authorize request: %s. Redirecting back to login page...", err.Error()))
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
		return
	}

	// clear cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: false,
	})

	// remove the session & CSRF tokens from the database
	username := r.FormValue("username")
	err = db.LogoutUser(username)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Failed to logout user: %s", err.Error()))
		logs.Logs(logWarning, "User session & CSRF tokens have not been removed from the database")
		http.Error(w, fmt.Sprintf("Unable to logout user: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	logs.Logs(logInfo, "User logged out successfully. Redirected to index page...")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
