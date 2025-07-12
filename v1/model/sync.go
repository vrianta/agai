package model

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	DatabaseHandler "github.com/vrianta/agai/v1/database"
	"github.com/vrianta/agai/v1/internal/config"
)

// Function to get the table topology and compare with the latest FieldTypes and generate a new SQL queryBuilder to alter the table
// This function will be used to update the table structure if there are any changes in the FieldTypes
func (m *meta) syncTableSchema() {
	schemaMap := make(map[string]schema, len(m.schemas))
	for _, s := range m.schemas {
		schemaMap[s.field] = s
	}

	FieldTypeset := make(FieldTypeset, len(m.FieldTypes))
	for _, f := range m.FieldTypes {
		FieldTypeset[f.name] = f
	}

	var pendingAddFields []*Field

	reader := bufio.NewReader(os.Stdin)

	for _, field := range m.FieldTypes {
		schema, exists := schemaMap[field.name]
		if !exists {
			if config.GetBuild() {
				fmt.Printf("Field '%s' not in DB. Add? (y/n): ", field.name)
				if input, _ := reader.ReadString('\n'); strings.TrimSpace(input) != "y" {
					fmt.Printf("[AddField] Skipped: %s\n", field.name)
					continue
				}
			}
			// m.addField(field)
			pendingAddFields = append(pendingAddFields, field)
			continue
		}

		filed_type, field_length := schema.parseSQLType()
		shouldChange := false
		reasons := []string{}

		if !field.Compare(filed_type) {
			reasons = append(reasons, fmt.Sprintf("type mismatch(old:%s,new:%s)", filed_type, field.Type.string()))
			shouldChange = true
		}
		if !(field_length == 1 && field.Length == 0) && field_length != field.Length {
			reasons = append(reasons, fmt.Sprintf("length mismatch(old:%d:new:%d)", field_length, field.Length))
			shouldChange = true
		}
		if schema.defaultVal.String != field.DefaultValue {
			reasons = append(reasons, "default mismatch")
			shouldChange = true
		}
		if schema.nullable == "YES" && !field.Nullable {
			reasons = append(reasons, "nullable mismatch")
			shouldChange = true
		}
		if schema.nullable == "NO" && field.Nullable {
			reasons = append(reasons, "nullable mismatch")
			shouldChange = true
		}
		if schema.extra == "auto_increment" && !field.AutoIncrement {
			reasons = append(reasons, "auto_increment mismatch")
			shouldChange = true
		}

		if shouldChange {
			if config.GetBuild() {
				fmt.Printf("Field '%s' requires update (%s). Proceed? (y/n): ", field.name, strings.Join(reasons, ", "))
				input, _ := reader.ReadString('\n')
				if strings.TrimSpace(input) != "y" {
					fmt.Printf("\n[Modify] Skipped update of: %s\n", field.name)
					continue
				}
			}
			m.modifyDBField(field)
		}

		// Index differences
		if schema.isunique != field.Index.Unique {
			if config.GetBuild() {
				fmt.Printf("UNIQUE index mismatch on '%s'. Sync? (y/n): ", field.name)
				input, _ := reader.ReadString('\n')
				if strings.TrimSpace(input) != "y" {
					fmt.Printf("[Index] Skipped UNIQUE sync on: %s\n", field.name)
				} else {
					m.syncUniqueIndex(field, &schema)
				}
			} else {
				m.syncUniqueIndex(field, &schema)
			}
		}

		if schema.isprimary != field.Index.PrimaryKey {
			if config.GetBuild() {
				fmt.Printf("PRIMARY KEY mismatch on '%s'. Sync? (y/n): ", field.name)
				input, _ := reader.ReadString('\n')
				if strings.TrimSpace(input) != "y" {
					fmt.Printf("[Index] Skipped PRIMARY KEY sync on: %s\n", field.name)
				} else {
					m.syncPrimaryKey(field, &schema)
				}
			} else {
				m.syncPrimaryKey(field, &schema)
			}
		}

		if schema.isindex != field.Index.Index {
			if config.GetBuild() {
				fmt.Printf("INDEX mismatch on '%s'. Sync? (y/n): ", field.name)
				input, _ := reader.ReadString('\n')
				if strings.TrimSpace(input) != "y" {
					fmt.Printf("[Index] Skipped INDEX sync on: %s\n", field.name)
				} else {
					m.syncIndex(field, &schema)
				}
			} else {
				m.syncIndex(field, &schema)
			}
		}
	}

	// Check for fields to delete
	for _, schema := range m.schemas {
		if _, exists := FieldTypeset[schema.field]; !exists {
			if config.GetBuild() {
				fmt.Printf("Field '%s' exists in DB but not in model. Delete? (y/n): ", schema.field)
				input, err := reader.ReadString('\n')
				if err == nil && strings.TrimSpace(input) == "y" {
					m.removeDBField(schema.field)
				} else {
					fmt.Printf("[Delete] Skipped: %s\n", schema.field)
				}
			} else {
				fmt.Printf("Do you want to delete %s (y/n): ", schema.field)
				input, err := reader.ReadString('\n')
				if err == nil && strings.TrimSpace(input) == "y" {
					m.removeDBField(schema.field)
				} else {
					fmt.Printf("[Delete] Skipped: %s\n", schema.field)
				}
			}
		}
	}

	// Add all pending new fields now going to add in the table scema
	for _, field := range pendingAddFields {
		m.addField(field)
	}

}

// SyncModelSchema loads the current structure of the associated database table,
// including column definitions and index metadata (primary, unique, and standard indexes),
// and stores it in the model's internal schema list (m.schemas).
//
// This is used to detect schema differences for migration, validation, or syncing purposes.
// If the table does not exist, the function will exit early without error.
//
// Note: This function panics on database errors and should be called only when
// database availability is guaranteed.
func (m *meta) SyncModelSchema() {
	// Get the active database connection
	databaseObj, err := DatabaseHandler.GetDatabase()
	if err != nil {
		panic("Error getting database: " + err.Error())
	}

	// Check if the table for this model actually exists in the database
	checkqueryBuilder := `SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?`
	var count int
	err = databaseObj.QueryRow(checkqueryBuilder, m.TableName).Scan(&count)
	if err != nil {
		panic("Error checking table existence: " + err.Error())
	}
	if count == 0 {
		// If table does not exist, log and exit
		fmt.Printf("Table '%s' does not exist.\n", m.TableName)
		return
	}

	// Query the structure of the existing table
	rows, err := databaseObj.Query("SHOW COLUMNS FROM `" + m.TableName + "`")
	if err != nil {
		panic("Error getting old table structure: " + err.Error())
	}
	defer rows.Close() // Ensure result rows are closed

	// Clear any previously cached schema info
	m.schemas = nil

	// Query template to get index information for a column
	indexqueryBuilder := `
	SELECT 
	column_name, 
	index_name,
	non_unique
	FROM information_schema.statistics
	WHERE table_schema = ?
	AND table_name = ?
	AND column_name = ?`

	// Iterate through each column of the table
	for rows.Next() {
		_scema := schema{}
		// Scan each column's structure into the _scema struct
		if err := rows.Scan(&_scema.field, &_scema.fieldType, &_scema.nullable, &_scema.key, &_scema.defaultVal, &_scema.extra); err != nil {
			panic("Error scanning row: " + err.Error())
		}

		// Query the index info for this column
		if idxRows, err := databaseObj.Query(indexqueryBuilder, config.GetDatabaseConfig().Database, m.TableName, _scema.field); err != nil {
			panic("Error getting index information: " + err.Error())
		} else {
			defer idxRows.Close()

			// Process each index row
			for idxRows.Next() {
				var columnName, indexName string
				var nonUnique int
				if err := idxRows.Scan(&columnName, &indexName, &nonUnique); err != nil {
					panic("Error scanning index row: " + err.Error())
				}

				// Check if it's a primary key
				if indexName == "PRIMARY" {
					_scema.isprimary = true
				} else {
					// Determine if it's a standard index or unique constraint based on naming convention
					suffix := strings.Split(indexName, "_")
					switch suffix[0] {
					case "idx":
						_scema.isindex = true
					case "unq":
						_scema.isunique = true
					}
				}
			}
		}

		// Add the parsed schema to the model's schema list
		m.schemas = append(m.schemas, _scema)
	}
}
