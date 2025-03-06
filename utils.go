package main

import (
	"crypto/rand"
	"encoding/base64"
	"log"

	"golang.org/x/crypto/bcrypt"
)

// HashedPassword hashes the given password using bcrypt with a cost of 10.
// It returns the hashed password as a string and any error encountered.
func HashedPassword(password string) (string, error) {
	// byte representation of the password string, password hashed 2^10 times
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

/*
can also do:

func HashedPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
*/

// checkPasswordHash compares a hashed password with its possible plaintext equivalent.
// It returns true if the password matches the hash, otherwise false.
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// generateSessionToken generates a cryptographically secure random session token of length
// bytes and encodes it in base64. The token is suitable for use as a session ID or CSRF
// token in a web application.
func generateToken(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		log.Fatalf("Error generating session token: %v", err.Error())
	}
	return base64.URLEncoding.EncodeToString(bytes)
}
