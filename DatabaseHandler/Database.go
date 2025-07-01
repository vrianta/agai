package DatabaseHandler

import (
	"database/sql"
	"fmt"

	"github.com/vrianta/Server/Config"
)

// Function to init the Database with the Database/sql object and store it in the program
func Init() error {

	if Config.GetDatabaseConfig().Host == "" {

		return nil
	}
	var err error
	if database, err = sql.Open(Config.GetDatabaseDriver(), Config.GetDSN()); err != nil {
		return err
	}

	Initialized = true
	return nil
}

func GetDatabase() (*sql.DB, error) {
	if !Initialized {
		return nil, fmt.Errorf("Database configuration is not set")
	}
	if database == nil {
		return nil, fmt.Errorf("Database is not initialized")
	}
	return database, nil
}
