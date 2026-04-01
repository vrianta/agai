package model

import "github.com/vrianta/agai/v1/utils"

var MigrationName string = Date.string() + "_" + Time.string() + "_" + "migration.sql"

// migrationName is the file name
func recordMigrationSchema(queries string) {
	// check if the migrationName file exists
	if !utils.FileExists(MigrationName) {
		// create the file
		utils.CreateFile(MigrationName)
	}
	// append the queries to the file
	utils.AppendToFile(MigrationName, queries)
}
