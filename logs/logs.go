package logs

import "log"

const (
	info   = "INFO: "
	warn   = "WARNING: "
	logErr = "ERROR: "
	db     = "DATABASE: "
	dbErr  = "DATABASE ERROR: "
)

var logChannel = make(chan string)

// ProcessLogs continuously listens for log messages on the logChannel
// and prints them to the standard output.
func ProcessLogs() {
	for logMessage := range logChannel {
		log.Println(logMessage)
	}
}

/*
Logs writes a log message to the logChannel with the appropriate prefix.

logType must be one of the following:

1: info

2: warn

3: logErr

4: db

5: dbErr
*/
func Logs(logType int, logMessage string) {
	var loggedMessage string
	switch logType {
	case 1:
		loggedMessage = info + logMessage
	case 2:
		loggedMessage = warn + logMessage
	case 3:
		loggedMessage = logErr + logMessage
	case 4:
		loggedMessage = db + logMessage
	case 5:
		loggedMessage = dbErr + logMessage
	}

	logChannel <- loggedMessage
}
