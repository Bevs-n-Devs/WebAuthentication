package utils

import (
	"crypto/rand"
	"encoding/base64"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func ValidateUser(user string, password string) bool {
	if user == "pythonAkoto" && password == "password123" {
		return true
	}

	return false
}

// HashedPassword generates a hashed password using bcrypt, hashing the password
// 2^10 times. The function returns a string representation of the hashed
// password and an error if hashing fails.
func HashedPassword(password string) (string, error) {
	// byte representation of the password string, password hashed 2^10 times
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

// CheckPasswordHash takes a password and hash string and checks if the hash
// matches the password. The function returns true if the hash matches the
// password and false otherwise.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateToken generates a cryptographically secure random token of the given
// length and returns it as a string. The token is suitable for use as a session
// token in a web application.
func GenerateToken(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		log.Fatalf("Error generating session token: %v", err.Error())
	}
	return base64.URLEncoding.EncodeToString(bytes)
}
