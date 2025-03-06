package main

import (
	"errors"
	"net/http"
)

var AuthError = errors.New("Unauthorized")

func Authorize(r *http.Request) error {
	username := r.FormValue("username")
	// grab user data from mock database
	user, ok := usersDB[username]
	if !ok {
		return AuthError
	}

	// get the session token from the cookie
	sessionToken, err := r.Cookie("session_token")
	// check if cookie is not empty AND matching database session token
	if err != nil || sessionToken.Value == "" || sessionToken.Value != user.SessionToken {
		return AuthError
	}

	// get the CSRF token from the headers
	csrf := r.Header.Get("X-CSRF-Token")
	// check if header is not empty AND matching database CSRF token
	if csrf == "" || csrf != user.CSRFToken {
		return AuthError
	}

	return nil
}
