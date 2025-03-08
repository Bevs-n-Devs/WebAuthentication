package handlers

import "html/template"

const (
	logInfo    = 1
	logWarning = 2
	logErr     = 3
	logDb      = 4
)

var (
	Templates *template.Template // global variable to hold HTML templates

	htmlTemplate = template.Must(template.ParseFiles("./templates/index.html"))
)
