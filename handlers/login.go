package handlers

import (
	"fmt"
	"net/http"

	"github.com/Bevs-n-Devs/WebAuthentication/logs"
)

func Login(w http.ResponseWriter, r *http.Request) {
	err := Templates.ExecuteTemplate(w, "login.html", nil)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Failed to execute template: %s", err.Error()))
		http.Error(w, fmt.Sprintf("Unable to load page: %s", err.Error()), http.StatusInternalServerError)
	}
}
