package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Bevs-n-Devs/WebAuthentication/db"
	"github.com/Bevs-n-Devs/WebAuthentication/logs"
)

const (
	logInfo = 1
	logErr  = 3
)

var ErrAuth = errors.New("unauthorized - user not authenticated")

/*
AuthorizeRequest validates the session and CSRF tokens for the given HTTP request.
It extracts the username from the form data, retrieves the session token from cookies,
and the CSRF token from headers. It then checks these tokens against the database
for validity. If any token is missing or invalid, it returns an error indicating
unauthorized access. Returns nil if both tokens are valid.

Returns:

- error: An error if the session or CSRF tokens are missing or invalid.
*/
func AuthorizeRequest(r *http.Request) error {
	// get the session token from the cookie
	sessionToken, err := r.Cookie("session_token")
	if err != nil || sessionToken.Value == "" {
		return fmt.Errorf("%s! Session token is missing: %s", ErrAuth, err.Error())
	}

	username, err := db.GetUsernameFromSessionToken(sessionToken.Value)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Failed to get username from session token: %s", err.Error()))
		return err
	}

	// check if the username and session token are valid - a little redundant but good to have
	ok, err := db.ValidateSessionToken(username, sessionToken.Value)
	if err != nil {
		logs.Logs(logErr, fmt.Sprintf("Failed to validate session token: %s", err.Error()))
		return err
	}
	if !ok {
		logs.Logs(logErr, fmt.Sprintf("Invalid session token: %s", sessionToken.Value))
		return err
	}
	logs.Logs(logInfo, fmt.Sprintf("Session validation result: %t", ok))

	// get CSRF token from the cookie
	csrf, err := r.Cookie("csrf_token")
	if err != nil {
		return fmt.Errorf("%s! CSRF token is missing", ErrAuth)
	}

	// check if the username and CSRF token are valid
	ok, err = db.ValidateCSRFToken(username, csrf.Value)
	if err != nil {
		return fmt.Errorf("%s! Failed to validate CSRF token: %s", ErrAuth, err.Error())
	}
	if !ok {
		return fmt.Errorf("%s! Invalid CSRF token", ErrAuth)
	}
	logs.Logs(logInfo, fmt.Sprintf("CSRF validation result: %t", ok))

	return nil
}
