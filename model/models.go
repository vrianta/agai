package model

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	config "github.com/vrianta/Server/config"
	DatabaseHandler "github.com/vrianta/Server/database"
)

/*
 * This Package is to handle model in the database checking and creating tables and providing default functions to handle them
 * It will create the table,
 * It will update the table accordingly during the initial program startup only if the build is not true
 * So Dynaimic Table Updation will be handled during development only
 * It will provide the default functions to handle the model like Create, Read, Update, Delete
 */

func New(tableName string, fields map[string]Field) *Struct {
	_model := Struct{
		TableName: tableName,
		fields:    fields,
		primary: func(fields map[string]Field) *Field {
			for _, fields_val := range fields {
				if fields_val.Index.PrimaryKey {
					return &fields_val
				}
			}
			return nil
		}(fields),
	}

	_model.validate()

	ModelsRegistry[tableName] = &_model
	return &_model
}

func (m *Struct) validate() {
	primaryKeyCount := 0
	fieldNames := make(map[string]struct{})

	for _, field := range m.fields {
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
func (m *Struct) GetTableScema() {
	databaseObj, err := DatabaseHandler.GetDatabase()
	if err != nil {
		panic("Error getting database: " + err.Error())
	}

	// 1. Load column info
	checkQuery := `SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?`
	var count int
	err = databaseObj.QueryRow(checkQuery, m.TableName).Scan(&count)
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

	indexQuery := `
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

		if idxRows, err := databaseObj.Query(indexQuery, config.GetDatabaseConfig().Database, m.TableName, _scema.field); err != nil {
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

// Function to get the table topology and compare with the latest fields and generate a new SQL query to alter the table
// This function will be used to update the table structure if there are any changes in the fields
func (m *Struct) syncTableSchema() {
	// Use maps for faster lookups
	schemaMap := make(map[string]schema, len(m.schemas))
	for _, s := range m.schemas {
		schemaMap[s.field] = s
	}

	fieldMap := make(map[string]Field, len(m.fields))
	for _, f := range m.fields {
		fieldMap[f.Name] = f
	}

	// Check for new or changed fields
	for _, field := range m.fields {
		schema, exists := schemaMap[field.Name]
		if !exists {
			m.addField(&field)
			continue
		}

		filed_type, field_length := schema.parseSQLType()
		shouldChange := false

		// fmt.Println(schema)

		if filed_type != field.Type.string() {
			fmt.Printf("[Type Mismatch] Field: %-20s | DB: %-10s | Model: %-10s\n", field.Name, filed_type, field.Type.string())
			shouldChange = true
		}
		if !(field_length == 1 && field.Length == 0) && field_length != field.Length {
			fmt.Printf("[Length Mismatch] Field: %-20s | DB: %-5d | Model: %-5d\n", field.Name, field_length, field.Length)
			shouldChange = true
		}
		if schema.defaultVal.String != field.DefaultValue {
			fmt.Printf("[Default Value]  Field: %-20s | DB: %-10s | Model: %-10s\n", field.Name, schema.defaultVal.String, field.DefaultValue)
			shouldChange = true
		}
		if schema.nullable == "YES" && !field.Nullable {
			fmt.Printf("[Nullable]      Field: %-20s | DB: YES        | Model: NOT NULL\n", field.Name)
			shouldChange = true
		}
		if schema.nullable == "NO" && field.Nullable {
			fmt.Printf("[Nullable]      Field: %-20s | DB: NO         | Model: NULL\n", field.Name)
			shouldChange = true
		}
		if schema.extra == "auto_increment" && !field.AutoIncrement {
			fmt.Printf("[AutoIncrement] Field: %-20s | DB: auto_increment | Model: not auto_increment\n", field.Name)
			shouldChange = true
		}

		if shouldChange {
			m.updateField(&field)
		}

		// Check for index mismatches
		if schema.isunique != field.Index.Unique {
			// fmt.Println("unique are different")
			m.updateUniqueIndex(&field, &schema)
		}
		if schema.isprimary != field.Index.PrimaryKey {
			fmt.Println("Primary Keys are different")
			m.updatePrimaryKey(&field, &schema)
		}
		if schema.isindex != field.Index.Index {
			m.updateNormalIndex(&field, &schema)
		}
	}

	// Check for fields to delete
	for _, schema := range m.schemas {
		if _, exists := fieldMap[schema.field]; !exists {
			fmt.Printf("Do you want to delete %s (y/n): ", schema.field)
			reader := bufio.NewReader(os.Stdin)
			if input, err := reader.ReadString('\n'); err == nil && strings.TrimSpace(input) == "y" {
				m.dropField(schema.field)
			} else if err != nil {
				fmt.Printf("Error Getting Input: %s", err.Error())
			} else {
				fmt.Printf("[Delete]        Skipping deletion of field: %-20s\n", schema.field)
			}
		}
	}
}

func (m *Struct) CreateTableIfNotExists() {
	if len(m.schemas) > 0 { // if the lenth is more that 0 that means talbe is already created and no need to create it again instead we should focus on updating it
		// if !Config.GetBuild() { // table syncing will only work only if it is a build version
		m.syncTableSchema()
		// }
		return
	}
	sql := "CREATE TABLE IF NOT EXISTS " + m.TableName + " (\n"
	fieldDefs := []string{}

	for _, field := range m.fields {
		fieldDefs = append(fieldDefs, field.String())
	}

	sql += strings.Join(fieldDefs, ",\n")
	sql += "\n);"

	fmt.Println("\n[SQL] Table Creation Statement:\n" + sql + "\n")

	databaseObj, err := DatabaseHandler.GetDatabase()
	if err != nil {
		panic("Error getting database: " + err.Error())
	}

	_, err = databaseObj.Exec(sql)
	if err != nil {
		panic("Error creating table: " + err.Error() + "\nQuery:" + sql)
	}

	fmt.Printf("[Success] Table created or already exists: %s\n", m.TableName)
}

func (m *Struct) LoadIndexes() {
	db, err := DatabaseHandler.GetDatabase()
	if err != nil {
		panic(err)
	}

	query := `
	SELECT column_name, index_name, non_unique
	FROM information_schema.statistics
	WHERE table_schema = DATABASE() AND table_name = ?`

	rows, err := db.Query(query, m.TableName)
	if err != nil {
		panic("Error fetching index info: " + err.Error())
	}
	defer rows.Close()

	m.indexes = make(map[string]indexInfo)

	for rows.Next() {
		var col, idx string
		var nonUnique int
		if err := rows.Scan(&col, &idx, &nonUnique); err != nil {
			panic("Error scanning index row: " + err.Error())
		}
		m.indexes[col] = indexInfo{
			ColumnName: col,
			IndexName:  idx,
			NonUnique:  nonUnique == 1,
		}
	}
}

// function to add the new field in the table
func (m *Struct) addField(field *Field) {
	/*
		ALTER TABLE `users`
		ADD `newel` VARCHAR(20) NULL DEFAULT 'dwads' AFTER `userId`,
		ADD INDEX (`newel`);
	*/

	response := "ALTER TABLE `" + m.TableName + "`\n"
	response += "ADD " + field.columnDefinition() + field.addIndexStatement() + ";"

	if databaseObj, err := DatabaseHandler.GetDatabase(); err != nil {
		panic("Error While Adding new Field to the table" + err.Error())
	} else {
		if _, sql_err := databaseObj.Exec(response); sql_err != nil {
			panic("Error While Updating the Table Field" + sql_err.Error())
		} else {
			fmt.Printf("[AddField]      Table: %-20s | Field Added: %-20s\n", m.TableName, field.Name)
		}
	}
}

// function to change the field details
func (m *Struct) updateField(field *Field) {
	// ALTER TABLE `users` CHANGE `userId` `userId` INT(30) NOT NULL AUTO_INCREMENT;
	response := "ALTER TABLE `" + m.TableName + "` "
	response += "CHANGE `" + field.Name + "` " + field.columnDefinition() + ";"

	if databaseObj, err := DatabaseHandler.GetDatabase(); err != nil {
		panic("Error While Changing Field" + err.Error())
	} else {
		if _, sql_err := databaseObj.Exec(response); sql_err != nil {
			panic("Error While Changing the Table Field" + sql_err.Error() + "SQL QUERY: " + response)
		} else {
			fmt.Printf("[UpdateField]   Table: %-20s | Field Updated: %-20s\n", m.TableName, field.Name)
		}
	}
}

// Handles adding/dropping PRIMARY KEY
func (m *Struct) updatePrimaryKey(field *Field, schema *schema) {
	databaseObj, err := DatabaseHandler.GetDatabase()
	if err != nil {
		fmt.Println("Error updating primary key:", err)
		return
	}
	if schema.isprimary && !field.Index.PrimaryKey {
		// Drop primary key
		query := fmt.Sprintf("ALTER TABLE `%s` DROP PRIMARY KEY;", m.TableName)
		if _, err := databaseObj.Exec(query); err != nil {
			fmt.Println("[Index] Error dropping PRIMARY KEY:", err)
		} else {
			fmt.Printf("[Index] PRIMARY KEY dropped for field: %s\n", field.Name)
		}
	}
	if !schema.isprimary && field.Index.PrimaryKey {
		// Add primary key
		query := "ALTER TABLE " + m.TableName + " ADD PRIMARY KEY (" + field.Name + ")"
		if _, err := databaseObj.Query(query); err != nil {
			fmt.Println("[ERROR] failed to Add Primary Key ", err.Error())
			fmt.Println("[FAILED] Failed Query to Update Primary Key is: ", query)
		}
	}
}

// Handles adding/dropping UNIQUE index
func (m *Struct) updateUniqueIndex(field *Field, schema *schema) {
	databaseObj, err := DatabaseHandler.GetDatabase()
	if err != nil {
		fmt.Println("Error updating unique index:", err)
		return
	}
	indexName := fmt.Sprintf("unq_%s", field.Name)
	if schema.isunique && !field.Index.Unique {
		// Drop unique index
		query := fmt.Sprintf("ALTER TABLE `%s` DROP INDEX `%s`;", m.TableName, indexName)
		if _, err := databaseObj.Exec(query); err != nil {
			fmt.Println("[Index] Error dropping UNIQUE:", err)
		} else {
			fmt.Printf("[Index] UNIQUE dropped for field: %s\n", field.Name)
		}
	}
	if !schema.isunique && field.Index.Unique {
		// Add unique index
		query := fmt.Sprintf("ALTER TABLE `%s` ADD UNIQUE `%s` (`%s`);", m.TableName, indexName, field.Name)
		if _, err := databaseObj.Exec(query); err != nil {
			fmt.Println("[Index] Error adding UNIQUE:", err)
		} else {
			fmt.Printf("[Index] UNIQUE added for field: %s\n", field.Name)
		}
	}
}

// Handles adding/dropping normal INDEX
func (m *Struct) updateNormalIndex(field *Field, schema *schema) {
	databaseObj, err := DatabaseHandler.GetDatabase()
	if err != nil {
		fmt.Println("Error updating index:", err)
		return
	}
	indexName := fmt.Sprintf("idx_%s", field.Name)
	if schema.isindex && !field.Index.Index {
		// Drop index
		query := fmt.Sprintf("ALTER TABLE `%s` DROP INDEX `%s`;", m.TableName, indexName)
		if _, err := databaseObj.Exec(query); err != nil {
			fmt.Println("[Index] Error dropping INDEX:", err)
		} else {
			fmt.Printf("[Index] INDEX dropped for field: %s\n", field.Name)
		}
	}
	if !schema.isindex && field.Index.Index {
		// Add index
		query := fmt.Sprintf("ALTER TABLE `%s` ADD INDEX `%s` (`%s`);", m.TableName, indexName, field.Name)
		if _, err := databaseObj.Exec(query); err != nil {
			fmt.Println("[Index] Error adding INDEX:", err)
		} else {
			fmt.Printf("[Index] INDEX added for field: %s\n", field.Name)
		}
	}
}

// Drop a field from the databasess
func (m *Struct) dropField(fieldName string) {
	//ALTER TABLE `users` DROP `userId`;
	query := "ALTER TABLE `" + m.TableName + "` DROP `" + fieldName + "`;"
	if databaseObj, err := DatabaseHandler.GetDatabase(); err != nil {
		panic("Error While Deleting Field" + err.Error())
	} else {
		if _, sql_err := databaseObj.Exec(query); sql_err != nil {
			panic("Error While Deleting the Field" + sql_err.Error())
		} else {
			fmt.Printf("[DropField]     Table: %-20s | Field Dropped: %-20s\n", m.TableName, fieldName)
		}
	}
}

// get the table name
func (m *Struct) GetTableName() string {
	return m.TableName
}

// function to get data of a field
func (m *Struct) GetFieldValue(field_name string) (any, error) {
	if field, ok := m.fields[field_name]; !ok {
		return nil, errors.New("")
	} else {
		return field.value, nil
	}
}

// Convert the Fetched Data to a of objects
// This function will convert the Struct to a map[string]any for easy access and manipulation
func (m *Struct) ToMap() map[string]any {
	response := make(map[string]any, len(m.fields))
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, field := range m.fields {
		wg.Add(1)

		go func(f Field) {
			defer wg.Done()

			var value any
			if f.value != nil {
				value = f.value
			} else {
				value = nil
			}

			mu.Lock()
			response[f.Name] = value
			mu.Unlock()
		}(field)
	}

	wg.Wait()
	return response
}

// Insert inserts a new record into the table using the provided values map.
// This is a dedicated Create/Insert function that does not overlap with table creation or schema management.
func (m *Struct) Insert(values map[string]any) error {
	q := m.Create()
	for k, v := range values {
		q.insertFields[k] = v
	}
	return q.Exec()
}

func (m *Struct) GetPrimary() *Field {
	if !m.PrimaryKeyExists() {
		panic("You asked for Primary Key but the Model do not have any PrimaryKey")
	}
	return m.primary
}
func (m *Struct) GetPrimaryKeyVal() string {
	return m.primary.value.(string)
}

/*
To check if the model has primary key or not

true ->  if exists
false -> if not exists
*/
func (m *Struct) PrimaryKeyExists() bool {
	return m.primary != nil
}
