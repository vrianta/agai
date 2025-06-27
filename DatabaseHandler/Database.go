package DatabaseHandler

import (
	"database/sql"
	"fmt"
)

// Function to init the Database with the Database/sql object and store it in the program
func Init(_sql *sql.DB) error {
	if _sql == nil {
		return fmt.Errorf("Database Initialisation failed: provided Database Object is nil")
	}
	// Store the database object in a global variable or a struct
	// This is just an example, you can modify it as per your program structure
	database = _sql
	return nil
}

func GetDatabase() (*sql.DB, error) {
	if database == nil {
		return nil, fmt.Errorf("Database is not initialized")
	}
	return database, nil
}
