package handlers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Bevs-n-Devs/WebAuthentication/logs"
)

func StartHTTPServer() {
	logs.Logs(logInfo, "Starting HTTP server...")

	// initialize templates
	InitTemplates()

	// static file server for assets like CSS (if any)
	// static directory needed in project root
	var staticFiles = http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", staticFiles))

	// define routes
	http.HandleFunc("/", IndexRoute)
	http.HandleFunc("/account", Account)
	http.HandleFunc("/create-account", CreateAccount)
	http.HandleFunc("/login", Login)
	http.HandleFunc("/submit-login", SubmitLogin)
	http.HandleFunc("/dashboard", Dashboard)
	http.HandleFunc("/logout", LogoutUser)

	// initialize port
	httpPort := os.Getenv("PORT")
	// start server on local machine
	if httpPort == "" {
		logs.Logs(logWarning, "Could not get PORT from hosting platform. Deafaulting to http://localhost:9003...")
		httpPort = "9003"
		err := http.ListenAndServe(fmt.Sprintf(":%s", httpPort), nil)
		if err != nil {
			logs.Logs(logErr, fmt.Sprintf("Failed to start HTTP server: %s", err.Error()))
		}
	}

	// start the server on hosting platform
	logs.Logs(logInfo, fmt.Sprintf("HTTP server started on http://localhost:%s", httpPort))
	err := http.ListenAndServe(fmt.Sprintf(":%s", httpPort), nil)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Failed to start HTTP server: %s", err.Error()))
	}
}
