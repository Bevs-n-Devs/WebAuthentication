package handlers

import (
	"fmt"
	"net/http"

	"github.com/Bevs-n-Devs/WebAuthentication/logs"
)

func Account(w http.ResponseWriter, r *http.Request) {
	err := Templates.ExecuteTemplate(w, "account.html", nil)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Failed to execute template: %s", err.Error()))
		http.Error(w, fmt.Sprintf("Unable to load page: %s", err.Error()), http.StatusInternalServerError)
	}
}
