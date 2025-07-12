package database

import "database/sql"

var (
	database    *sql.DB // Global variable to hold the database connection
	Initialized = false
)
