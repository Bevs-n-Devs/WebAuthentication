package db

import "database/sql"

const (
	logWarning = 2
	logErr     = 3
	logDb      = 4
	logDbErr   = 5
)

var (
	db *sql.DB // global DB variable to hold DB connection
)
