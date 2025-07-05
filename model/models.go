package model

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	DatabaseHandler "github.com/vrianta/Server/database"
)

/*
 * This Package is to handle model in the database checking and creating tables and providing default functions to handle them
 * It will create the table,
 * It will update the table accordingly during the initial program startup only if the build is not true
 * So Dynaimic Table Updation will be handled during development only
 * It will provide the default functions to handle the model like Create, Read, Update, Delete
 */

func newModel(tableName string, FieldTypes FieldTypeset) Table {
	_model := Table{
		TableName:  tableName,
		FieldTypes: FieldTypes,
		primary: func(FieldTypes FieldTypeset) *Field {
			for _, field := range FieldTypes {
				if field.Index.PrimaryKey {
					return field // Return the pointer directly from the map
				}
			}
			return nil
		}(FieldTypes),
	}

	_model.loadSchemaFromDB()

	return _model
}

func New[T any](tableName string, structure T) struct {
	Table
	Definition T
} {
	t := reflect.TypeOf(structure)
	v := reflect.ValueOf(structure)

	if t.Kind() != reflect.Struct {
		panic("structure passed to New must be a struct")
	}

	FieldTypeset := make(FieldTypeset, t.NumField())

	for i := 0; i < t.NumField(); i++ {
		structField := t.Field(i) // get metadata (e.g. "Element", "Value")
		value := v.Field(i).Interface()

		field, ok := value.(Field)
		if !ok {
			panic(fmt.Sprintf("[Model Error] Field '%s' is not of type model.Field", structField.Name))
		}

		// Always override the Name field based on struct variable name
		field.name = structField.Name

		FieldTypeset[structField.Name] = &field
	}

	response := struct {
		Table
		Definition T
	}{
		Table:      newModel(tableName, FieldTypeset),
		Definition: structure,
	}

	ModelsRegistry[tableName] = &response.Table
	return response
}

func (m *Table) CreateTableIfNotExists() {
	if len(m.schemas) > 0 { // if the lenth is more that 0 that means talbe is already created and no need to create it again instead we should focus on updating it
		// if !Config.GetBuild() { // table syncing will only work only if it is a build version
		m.syncTableSchema()
		// }
		return
	}
	sql := "CREATE TABLE IF NOT EXISTS " + m.TableName + " (\n"
	fieldDefs := []string{}

	for _, field := range m.FieldTypes {
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
		panic("Error creating table: " + err.Error() + "\nqueryBuilder:" + sql)
	}

	fmt.Printf("[Success] Table created or already exists: %s\n", m.TableName)
}

// func (m *Table) loadIndexMetadata() {
// 	db, err := DatabaseHandler.GetDatabase()
// 	if err != nil {
// 		panic(err)
// 	}

// 	queryBuilder := `
// 	SELECT column_name, index_name, non_unique
// 	FROM information_schema.statistics
// 	WHERE table_schema = DATABASE() AND table_name = ?`

// 	rows, err := db.Query(queryBuilder, m.TableName)
// 	if err != nil {
// 		panic("Error fetching index info: " + err.Error())
// 	}
// 	defer rows.Close()

// 	m.indexes = make(map[string]indexInfo)

// 	for rows.Next() {
// 		var col, idx string
// 		var nonUnique int
// 		if err := rows.Scan(&col, &idx, &nonUnique); err != nil {
// 			panic("Error scanning index row: " + err.Error())
// 		}
// 		m.indexes[col] = indexInfo{
// 			ColumnName: col,
// 			IndexName:  idx,
// 			NonUnique:  nonUnique == 1,
// 		}
// 	}
// }

// Handles adding/dropping PRIMARY KEY
func (m *Table) syncPrimaryKey(field *Field, schema *schema) {
	databaseObj, err := DatabaseHandler.GetDatabase()
	if err != nil {
		fmt.Println("Error updating primary key:", err)
		return
	}
	if schema.isprimary && !field.Index.PrimaryKey {
		// Drop primary key
		queryBuilder := fmt.Sprintf("ALTER TABLE `%s` DROP PRIMARY KEY;", m.TableName)
		if _, err := databaseObj.Exec(queryBuilder); err != nil {
			fmt.Println("[Index] Error dropping PRIMARY KEY:", err)
		} else {
			fmt.Printf("[Index] PRIMARY KEY dropped for field: %s\n", field.name)
		}
	}
	if !schema.isprimary && field.Index.PrimaryKey {
		// Add primary key
		queryBuilder := "ALTER TABLE " + m.TableName + " ADD PRIMARY KEY (" + field.name + ")"
		if _, err := databaseObj.Query(queryBuilder); err != nil {
			fmt.Println("[ERROR] failed to Add Primary Key ", err.Error())
			fmt.Println("[FAILED] Failed queryBuilder to Update Primary Key is: ", queryBuilder)
		}
	}
}

// Handles adding/dropping UNIQUE index
func (m *Table) syncUniqueIndex(field *Field, schema *schema) {
	databaseObj, err := DatabaseHandler.GetDatabase()
	if err != nil {
		fmt.Println("Error updating unique index:", err)
		return
	}
	indexName := fmt.Sprintf("unq_%s", field.name)
	if schema.isunique && !field.Index.Unique {
		// Drop unique index
		queryBuilder := fmt.Sprintf("ALTER TABLE `%s` DROP INDEX `%s`;", m.TableName, indexName)
		if _, err := databaseObj.Exec(queryBuilder); err != nil {
			fmt.Println("[Index] Error dropping UNIQUE:", err)
		} else {
			fmt.Printf("[Index] UNIQUE dropped for field: %s\n", field.name)
		}
	}
	if !schema.isunique && field.Index.Unique {
		// Add unique index
		queryBuilder := fmt.Sprintf("ALTER TABLE `%s` ADD UNIQUE `%s` (`%s`);", m.TableName, indexName, field.name)
		if _, err := databaseObj.Exec(queryBuilder); err != nil {
			fmt.Println("[Index] Error adding UNIQUE:", err)
		} else {
			fmt.Printf("[Index] UNIQUE added for field: %s\n", field.name)
		}
	}
}

// Handles adding/dropping normal INDEX
func (m *Table) syncIndex(field *Field, schema *schema) {
	databaseObj, err := DatabaseHandler.GetDatabase()
	if err != nil {
		fmt.Println("Error updating index:", err)
		return
	}
	indexName := fmt.Sprintf("idx_%s", field.name)
	if schema.isindex && !field.Index.Index {
		// Drop index
		queryBuilder := fmt.Sprintf("ALTER TABLE `%s` DROP INDEX `%s`;", m.TableName, indexName)
		if _, err := databaseObj.Exec(queryBuilder); err != nil {
			fmt.Println("[Index] Error dropping INDEX:", err)
		} else {
			fmt.Printf("[Index] INDEX dropped for field: %s\n", field.name)
		}
	}
	if !schema.isindex && field.Index.Index {
		// Add index
		queryBuilder := fmt.Sprintf("ALTER TABLE `%s` ADD INDEX `%s` (`%s`);", m.TableName, indexName, field.name)
		if _, err := databaseObj.Exec(queryBuilder); err != nil {
			fmt.Println("[Index] Error adding INDEX:", err)
		} else {
			fmt.Printf("[Index] INDEX added for field: %s\n", field.name)
		}
	}
}

// get the table name
func (m *Table) GetTableName() string {
	return m.TableName
}

// Convert the Fetched Data to a of objects
// This function will convert the Table to a map[string]any for easy access and manipulation
// func (m *Table) ToMap() map[string]any {
// 	response := make(map[string]any, len(m.FieldTypes))
// 	var wg sync.WaitGroup
// 	var mu sync.Mutex

// 	for _, field := range m.FieldTypes {
// 		wg.Add(1)

// 		go func(f *Field) {
// 			defer wg.Done()

// 			var value any
// 			if f.value != nil {
// 				value = f.value
// 			} else {
// 				value = nil
// 			}

// 			mu.Lock()
// 			response[f.Name] = value
// 			mu.Unlock()
// 		}(&field)
// 	}

// 	wg.Wait()
// 	return response
// }

// InsertRow InsertRows a new record into the table using the provided values map.
// This is a dedicated Create/InsertRow function that does not overlap with table creation or schema management.
func (m *Table) InsertRow(values map[string]any) error {
	q := m.Create()
	for k, v := range values {
		q.InsertRowFieldTypes[k] = v
	}
	return q.Exec()
}

func (m *Table) GetPrimaryKey() *Field {
	if !m.HasPrimaryKey() {
		panic("Primary Key is Required for but the Model(" + m.TableName + ") ")
	}
	return m.primary
}

/*
To check if the model has primary key or not

true ->  if exists
false -> if not exists
*/
func (m *Table) HasPrimaryKey() bool {
	if m.primary != nil {
		return true
	}
	for _, field := range m.FieldTypes {
		if field.Index.PrimaryKey {
			m.primary = field
			return true // Return the pointer directly from the map
		}
	}

	return false
}

/*
GetField(fieldname) -> return pointer of the field
*/

func (m *Table) GetField(field_name string) *Field {
	field, ok := m.FieldTypes[field_name]
	if !ok {
		return nil
	}
	return field
}

/*
GetField(fieldname) -> return pointer of the field
*/

func (m *Table) GetFieldTypes() *FieldTypeset {
	return &m.FieldTypes
}

// Print the Objects of the models as the good for debug perpose
func (r *Results) PrintAsTable() {
	if len(*r) == 0 {
		return
	}

	// Collect all unique column names across all rows
	colSet := map[string]struct{}{}
	for _, row := range *r {
		for col := range row {
			colSet[col] = struct{}{}
		}
	}

	// Sort column names for consistent display
	var colNames []string
	for col := range colSet {
		colNames = append(colNames, col)
	}
	sort.Strings(colNames)

	// Print header
	for _, col := range colNames {
		fmt.Printf("| %-15s", col)
	}
	fmt.Println("|")

	// Print separator
	fmt.Println(strings.Repeat("-", len(colNames)*18))

	// Print each row
	for _, row := range *r {
		for _, col := range colNames {
			val := row[col]
			fmt.Printf("| %-15v", val)
		}
		fmt.Println("|")
	}
}
