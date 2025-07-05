package model

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/vrianta/Server/config"
	DatabaseHandler "github.com/vrianta/Server/database"
)

// Function to get the table topology and compare with the latest FieldTypes and generate a new SQL queryBuilder to alter the table
// This function will be used to update the table structure if there are any changes in the FieldTypes
func (m *Table) syncTableSchema() {
	schemaMap := make(map[string]schema, len(m.schemas))
	for _, s := range m.schemas {
		schemaMap[s.field] = s
	}

	FieldTypeset := make(FieldTypeset, len(m.FieldTypes))
	for _, f := range m.FieldTypes {
		FieldTypeset[f.Name] = f
	}

	reader := bufio.NewReader(os.Stdin)

	for _, field := range m.FieldTypes {
		schema, exists := schemaMap[field.Name]
		if !exists {
			if config.GetBuild() {
				fmt.Printf("Field '%s' not in DB. Add? (y/n): ", field.Name)
				if input, _ := reader.ReadString('\n'); strings.TrimSpace(input) != "y" {
					fmt.Printf("[AddField] Skipped: %s\n", field.Name)
					continue
				}
			}
			m.addField(field)
			continue
		}

		filed_type, field_length := schema.parseSQLType()
		shouldChange := false
		reasons := []string{}

		if filed_type != field.Type.string() {
			reasons = append(reasons, "type mismatch")
			shouldChange = true
		}
		if !(field_length == 1 && field.Length == 0) && field_length != field.Length {
			reasons = append(reasons, "length mismatch")
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
				fmt.Printf("Field '%s' requires update (%s). Proceed? (y/n): ", field.Name, strings.Join(reasons, ", "))
				input, _ := reader.ReadString('\n')
				if strings.TrimSpace(input) != "y" {
					fmt.Printf("[Modify] Skipped update of: %s\n", field.Name)
					continue
				}
			}
			m.modifyDBField(field)
		}

		// Index differences
		if schema.isunique != field.Index.Unique {
			if config.GetBuild() {
				fmt.Printf("UNIQUE index mismatch on '%s'. Sync? (y/n): ", field.Name)
				input, _ := reader.ReadString('\n')
				if strings.TrimSpace(input) != "y" {
					fmt.Printf("[Index] Skipped UNIQUE sync on: %s\n", field.Name)
				} else {
					m.syncUniqueIndex(field, &schema)
				}
			} else {
				m.syncUniqueIndex(field, &schema)
			}
		}

		if schema.isprimary != field.Index.PrimaryKey {
			if config.GetBuild() {
				fmt.Printf("PRIMARY KEY mismatch on '%s'. Sync? (y/n): ", field.Name)
				input, _ := reader.ReadString('\n')
				if strings.TrimSpace(input) != "y" {
					fmt.Printf("[Index] Skipped PRIMARY KEY sync on: %s\n", field.Name)
				} else {
					m.syncPrimaryKey(field, &schema)
				}
			} else {
				m.syncPrimaryKey(field, &schema)
			}
		}

		if schema.isindex != field.Index.Index {
			if config.GetBuild() {
				fmt.Printf("INDEX mismatch on '%s'. Sync? (y/n): ", field.Name)
				input, _ := reader.ReadString('\n')
				if strings.TrimSpace(input) != "y" {
					fmt.Printf("[Index] Skipped INDEX sync on: %s\n", field.Name)
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
}

func (m *Table) loadSchemaFromDB() {
	primaryKeyCount := 0
	fieldNames := make(map[string]struct{})

	for _, field := range m.FieldTypes {
		// Check for duplicate field names
		if _, exists := fieldNames[field.Name]; exists {
			panic(fmt.Sprintf("[Validation Error] Duplicate field name '%s' in Table '%s'.\n", field.Name, m.TableName))
		}
		fieldNames[field.Name] = struct{}{}

		// PRIMARY KEY and UNIQUE cannot both be true
		if field.Index.PrimaryKey && field.Index.Unique {
			panic(fmt.Sprintf("[Validation Error] Field '%s' in Table '%s' cannot be both PRIMARY KEY and UNIQUE.\n", field.Name, m.TableName))
		}

		// Count primary keys
		if field.Index.PrimaryKey {
			primaryKeyCount++
			// PRIMARY KEY must not be nullable
			if field.Nullable {
				panic(fmt.Sprintf("[Validation Error] Field '%s' in Table '%s' is PRIMARY KEY but marked as nullable.\n", field.Name, m.TableName))
			}
			// PRIMARY KEY should not have default value
			if field.DefaultValue != "" {
				panic(fmt.Sprintf("[Validation Error] Field '%s' in Table '%s' is PRIMARY KEY but has a default value.\n", field.Name, m.TableName))
			}
		}

		// AutoIncrement should only be on integer types and primary key
		if field.AutoIncrement {
			if !field.Index.PrimaryKey {
				panic(fmt.Sprintf("[Validation Error] Field '%s' in Table '%s' is AUTO_INCREMENT but not PRIMARY KEY.\n", field.Name, m.TableName))
			}
			// You may want to check for integer type here, e.g.:
			if !strings.HasPrefix(strings.ToLower(field.Type.string()), "int") {
				panic(fmt.Sprintf("[Validation Error] Field '%s' in Table '%s' is AUTO_INCREMENT but not of integer type.\n", field.Name, m.TableName))
			}
		}
	}

	// Only one primary key allowed
	if primaryKeyCount > 1 {
		panic(fmt.Sprintf("[Validation Error] Table '%s' has more than one PRIMARY KEY field.\n", m.TableName))
	}
}

// Function to get the table scema of the mdoels and store them in the object
func (m *Table) SyncModelSchema() {
	databaseObj, err := DatabaseHandler.GetDatabase()
	if err != nil {
		panic("Error getting database: " + err.Error())
	}

	// 1. Load column info
	checkqueryBuilder := `SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?`
	var count int
	err = databaseObj.QueryRow(checkqueryBuilder, m.TableName).Scan(&count)
	if err != nil {
		panic("Error checking table existence: " + err.Error())
	}
	if count == 0 {
		fmt.Printf("Table '%s' does not exist.\n", m.TableName)
		return
	}

	rows, err := databaseObj.Query("SHOW COLUMNS FROM `" + m.TableName + "`")
	if err != nil {
		panic("Error getting old table structure: " + err.Error())
	}
	defer rows.Close()

	m.schemas = nil // clear any previous values

	indexqueryBuilder := `
	SELECT 
	column_name, 
	index_name,
	non_unique
	FROM information_schema.statistics
	WHERE table_schema = ?
	AND table_name = ?
	AND column_name = ?`

	for rows.Next() {
		_scema := schema{}
		if err := rows.Scan(&_scema.field, &_scema.fieldType, &_scema.nullable, &_scema.key, &_scema.defaultVal, &_scema.extra); err != nil {
			panic("Error scanning row: " + err.Error())
		}

		if idxRows, err := databaseObj.Query(indexqueryBuilder, config.GetDatabaseConfig().Database, m.TableName, _scema.field); err != nil {
			panic("Error getting index information: " + err.Error())
		} else {
			defer idxRows.Close()
			for idxRows.Next() {
				var columnName, indexName string
				var nonUnique int
				if err := idxRows.Scan(&columnName, &indexName, &nonUnique); err != nil {
					panic("Error scanning index row: " + err.Error())
				}

				if indexName == "PRIMARY" {
					_scema.isprimary = true
				} else {
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

		m.schemas = append(m.schemas, _scema)
	}

}
