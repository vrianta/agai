package database

import (
	"database/sql"
	"fmt"

	Config "github.com/vrianta/Server/config"
	log "github.com/vrianta/Server/log"
)

// Function to init the Database with the Database/sql object and store it in the program
func Init() {

	if Config.GetDatabaseConfig().Host == "" {
		log.WriteLog("DataBase Config do not have any host in it so We are skipping all database connections")
		return
	}
	var err error
	if database, err = sql.Open(Config.GetDatabaseDriver(), Config.GetDSN()); err != nil {
		panic("[ERROR] - DB Connection Failed Due to: " + err.Error())
	}

	Initialized = true
	fmt.Println("[Info] DataBase Connection Established Successfully")
}

func GetDatabase() (*sql.DB, error) {
	if !Initialized {
		panic("[ERROR] - You are calling for the Database but Connection with Database is not established properly")
	}
	if database == nil {
		panic("[ERROR] - You are calling for the Database but Connection with Database is not established properly")
	}
	return database, nil
}
