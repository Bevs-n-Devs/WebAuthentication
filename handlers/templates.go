package handlers

import (
	"fmt"
	"html/template"
	"os"

	"github.com/Bevs-n-Devs/WebAuthentication/logs"
)

func InitTemplates() {
	var err error
	Templates, err = template.ParseGlob("templates/*.html")
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Failed to parse templates: %s", err.Error()))
		os.Exit(1)
	}
}
