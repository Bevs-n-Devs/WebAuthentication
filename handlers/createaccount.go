package handlers

import (
	"fmt"
	"net/http"

	"github.com/Bevs-n-Devs/WebAuthentication/db"
	"github.com/Bevs-n-Devs/WebAuthentication/logs"
)

func CreateAccount(w http.ResponseWriter, r *http.Request) {
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

	err = db.CreateUser(username, password)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Failed to create user: %s", err.Error()))
		http.Error(w, fmt.Sprintf("Unable to load page: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	logs.Logs(logInfo, fmt.Sprintf("User %s created successfully. Redirected to login page...", username))
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
