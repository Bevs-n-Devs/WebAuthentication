package db

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	_ "embed"

	_ "github.com/lib/pq"

	"github.com/Bevs-n-Devs/WebAuthentication/env"
	"github.com/Bevs-n-Devs/WebAuthentication/logs"
	"github.com/Bevs-n-Devs/WebAuthentication/utils"
)

/*
ConnectDB connects to the PostgreSQL database via the DATABASE_URL environment
variable. If this variable is empty, it attempts to load the environment
variables from the .env file. The function logs the progress of the
connection attempt and returns an error if the connection cannot be
established.

Returns:

- error: An error object if the connection cannot be established.
*/
func ConnectDB() error {
	var err error

	// connect to database via environment variable
	if os.Getenv("DATABASE_URL") == "" {
		logs.Logs(logWarning, "Could not get database URL from hosting platform. Loading from .env file...")
		err := env.LoadEnv("env/.env")
		if err != nil {
			logs.Logs(logErr, fmt.Sprintf("Could not load environment variables from .env file: %s", err.Error()))
			return err
		}
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		logs.Logs(logDbErr, "Database URL is empty!")
		return fmt.Errorf("database URL is empty")
	}

	logs.Logs(logDb, "Connecting to database...")
	db, err = sql.Open("postgres", dbURL) // open db connection from global db variable
	if err != nil {
		logs.Logs(logDbErr, fmt.Sprintf("Could not connect to database: %s", err.Error()))
		return err
	}

	// verify connection
	logs.Logs(logDb, "Verifying database connection...")
	if db == nil {
		logs.Logs(logDbErr, "Database connection is empty!")
		return errors.New("database connection not established")
	}
	err = db.Ping()
	if err != nil {
		logs.Logs(logDbErr, fmt.Sprintf("Cannot ping database: %s", err.Error()))
		return err
	}
	logs.Logs(logDb, "Database connection established.")
	return nil
}

/*
CreateUser inserts a new user into the database with the provided username
and password. The password is hashed before being stored. It returns an
error if hashing the password fails or if the database execution encounters
an error.

Returns:

- error: An error if the database execution fails.
*/
func CreateUser(username, password string) error {
	if db == nil {
		logs.Logs(logDbErr, "Database connection is not initialized")
		return errors.New("database connection is not initialized")
	}

	hashedPwd, err := utils.HashedPassword(password)
	if err != nil {
		return err
	}

	query := `INSERT INTO tbl_web_auth_demo (username, hash_password) VALUES ($1, $2)`
	_, err = db.Exec(query, username, hashedPwd)
	return err
}

/*
UpdateSessionTokens generates new session and CSRF tokens for the given user
and updates them in the database. The tokens are valid for 24 hours. If the
update query fails, it returns an error. Otherwise, it returns the new
session token, CSRF token, and a nil error.

Returns:

- string: The new session token.

- string: The new CSRF token.

- error: An error if the update query fails.
*/
func UpdateSessionTokens(username string) (string, string, time.Time, error) {
	if db == nil {
		logs.Logs(logDbErr, "Database connection is not initialized")
		return "", "", time.Time{}, errors.New("database connection is not initialized")
	}

	sessionToken := utils.GenerateToken(32)
	csrfToken := utils.GenerateToken(32)
	expiry := time.Now().Add(24 * time.Hour) // 24-hour validity

	query := `UPDATE tbl_web_auth_demo SET session_token=$1, csrf_token=$2, token_expiry=$3 WHERE username=$4`
	_, err := db.Exec(query, sessionToken, csrfToken, expiry, username)
	if err != nil {
		logs.Logs(logDbErr, fmt.Sprintf("Failed to update session tokens: %s", err.Error()))
		return "", "", time.Time{}, err
	}

	logs.Logs(logDb, "Session tokens updated successfully")
	return sessionToken, csrfToken, expiry, nil
}

/*
AuthenticateUser checks if the provided username and password match the stored credentials.

It returns true if the credentials are correct, otherwise false. An error is returned if the query fails.

Returns:

- bool: True if the credentials are correct, otherwise false.

- error: An error if the query fails.
*/
func AuthenticateUser(username, password string) (bool, error) {
	if db == nil {
		logs.Logs(logDbErr, "Database connection is not initialized")
		return false, errors.New("database connection is not initialized")
	}

	var hashedPassword string
	query := `SELECT hash_password FROM tbl_web_auth_demo WHERE username=$1`
	err := db.QueryRow(query, username).Scan(&hashedPassword)
	if err != nil {
		return false, err
	}

	ok := utils.CheckPasswordHash(password, hashedPassword)
	if !ok {
		return false, errors.New("invalid password")
	}

	return true, nil
}

/*
ValidateSession checks if the given session token is valid for the given user.

It returns true if the session token is valid, otherwise false. An error is
returned if the query fails.

Returns:

- bool: True if the session token is valid, otherwise false.

- error: An error if the query fails.
*/
func ValidateSession(username, sessionToken string) (bool, error) {
	if db == nil {
		logs.Logs(logDbErr, "Database connection is not initialized")
		return false, errors.New("database connection is not initialized")
	}

	var expiry time.Time
	query := `SELECT token_expiry FROM tbl_web_auth_demo WHERE username=$1 AND session_token=$2`
	err := db.QueryRow(query, username, sessionToken).Scan(&expiry)
	if err != nil {
		return false, err
	}

	if time.Now().After(expiry) {
		return false, nil // Token expired
	}
	return true, nil
}

/*
ValidateSessionToken checks if the given session token is valid for the given user.
It returns true if the session token is valid, otherwise false. An error is
returned if the query fails.

Returns:

- bool: True if the session token is valid, otherwise false.

- error: An error if the query fails.
*/
func ValidateSessionToken(username, sessionToken string) (bool, error) {
	if db == nil {
		logs.Logs(logDbErr, "Database connection is not initialized")
		return false, errors.New("database connection is not initialized")
	}

	// query DB to get the stored session token
	var dbSessionToken string
	query := `
	SELECT session_token
	FROM tbl_web_auth_demo
	WHERE username = $1
	`
	err := db.QueryRow(query, username).Scan(&dbSessionToken)

	if err == sql.ErrNoRows {
		logs.Logs(logDbErr, "User not found")
		return false, errors.New("user not found")
	}

	if err != nil {
		logs.Logs(logDbErr, fmt.Sprintf("Failed to get session token: %s", err.Error()))
		return false, err
	}

	// compare the input session token with DB session token
	if sessionToken != dbSessionToken {
		logs.Logs(logDbErr, "Invalid session token")
		return false, nil
	}

	return true, nil
}

/*
ValidateCSRFToken checks if the given CSRF token is valid for the given user.
It returns true if the CSRF token is valid, otherwise false. An error is
returned if the query fails.

Returns:

- bool: True if the CSRF token is valid, otherwise false.

- error: An error if the query fails.
*/
func ValidateCSRFToken(username, csrfToken string) (bool, error) {
	if db == nil {
		logs.Logs(logDbErr, "Database connection is not initialized")
		return false, errors.New("database connection is not initialized")
	}

	// query DB to get the stored CSRF token
	var dbCSRFToken string
	query := `SELECT csrf_token FROM tbl_web_auth_demo WHERE username=$1`
	err := db.QueryRow(query, username).Scan(&dbCSRFToken)
	if err != nil {
		return false, err
	}

	// compare the input CSRF token with DB CSRF token
	if csrfToken != dbCSRFToken {
		return false, nil
	}
	return true, nil
}

/*
GetUsernameFromSessionToken retrieves the username associated with the given session token.

It returns the username if the session token is valid, otherwise an empty string. An error is
returned if the database query fails.

Returns:

- string: The username associated with the given session token.

- error: An error if the database query fails.
*/
func GetUsernameFromSessionToken(sessionToken string) (string, error) {
	if db == nil {
		logs.Logs(logDbErr, "Database connection is not initialized")
		return "", errors.New("database connection is not initialized")
	}

	var username string
	query := `SELECT username FROM tbl_web_auth_demo WHERE session_token=$1`
	err := db.QueryRow(query, sessionToken).Scan(&username)

	if err == sql.ErrNoRows {
		logs.Logs(logDbErr, "User not found")
		return "", errors.New("user not found")
	}

	if err != nil {
		logs.Logs(logDbErr, fmt.Sprintf("Failed to get session token: %s", err.Error()))
		return "", err
	}

	return username, nil
}

/*
LogoutUser removes the session and CSRF tokens for the given user by setting them to NULL,
effectively logging the user out. It also sets the token expiry to NULL. The function returns
an error if the database update query fails.

Returns:

- error: An error if the database update query fails.
*/
func LogoutUser(username string) error {
	if db == nil {
		logs.Logs(logDbErr, "Database connection is not initialized")
		return errors.New("database connection is not initialized")
	}

	query := `UPDATE tbl_web_auth_demo SET session_token=NULL, csrf_token=NULL, token_expiry=NULL WHERE username=$1`
	_, err := db.Exec(query, username)
	return err
}
