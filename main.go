package main

import (
	"fmt"
	"strings"

	"github.com/Bevs-n-Devs/WebAuthentication/db"
	"github.com/Bevs-n-Devs/WebAuthentication/handlers"
	"github.com/Bevs-n-Devs/WebAuthentication/logs"
)

const (
	logInfo  = 1
	logDbErr = 5
)

func main() {
	go logs.ProcessLogs()
	err := db.ConnectDB()
	if err != nil {
		logs.Logs(logDbErr, fmt.Sprintf("Failed to initialize database: %s", err.Error()))
	}

	go func() {
		handlers.StartHTTPServer()

		var templateNames []string
		for _, tmpl := range handlers.Templates.Templates() {
			templateNames = append(templateNames, tmpl.Name())
		}

		logs.Logs(logInfo, fmt.Sprintf("Loaded templates: %s", strings.Join(templateNames, ", ")))
	}()

	select {}
}
