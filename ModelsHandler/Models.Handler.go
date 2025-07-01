package Models

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/vrianta/Server/DatabaseHandler"
)

/*
 * This Package is to handle models in the database checking and creating tables and providing default functions to handle them
 * It will create the table,
 * It will update the table accordingly during the initial program startup only if the build is not true
 * So Dynaimic Table Updation will be handled during development only
 * It will provide the default functions to handle the models like Create, Read, Update, Delete
 */

func New(tableName string, fields []Field) *Struct {
	_model := Struct{
		TableName: tableName,
		Fields:    fields,
	}

	ModelsRegistry = append(ModelsRegistry, &_model)
	return &_model
}

// Function to get the table scema of the mdoels and store them in the object
func (m *Struct) GetTableScema() {
	databaseObj, err := DatabaseHandler.GetDatabase()

	if err != nil {
		panic("Error getting database: " + err.Error())
	}

	// Check if table exists
	checkQuery := `
		SELECT COUNT(*)
		FROM information_schema.tables
		WHERE table_schema = DATABASE()
		  AND table_name = ?`
	var count int
	err = databaseObj.QueryRow(checkQuery, m.TableName).Scan(&count)
	if err != nil {
		panic("Error checking table existence: " + err.Error())
	}
	if count == 0 {
		fmt.Printf("Table '%s' does not exist.\n", m.TableName)
		return // or handle gracefully
	}

	query := fmt.Sprintf("SHOW COLUMNS FROM `%s`", m.TableName)
	rows, err := databaseObj.Query(query)
	if err != nil {
		panic("Error getting old table structure: " + err.Error())
	}
	defer rows.Close()

	// Iterate over the rows (example)
	for rows.Next() {
		_scema := schema{}
		if err := rows.Scan(&_scema.field, &_scema.fieldType, &_scema.nullable, &_scema.key, &_scema.defaultVal, &_scema.extra); err != nil {
			panic("Error scanning row: " + err.Error())
		}

		m.schemas = append(m.schemas, _scema)
	}
}

// Function to get the table topology and compare with the latest fields and generate a new SQL query to alter the table
// This function will be used to update the table structure if there are any changes in the fields
func (m *Struct) SyncTableSchema() {
	// Use maps for faster lookups
	schemaMap := make(map[string]schema, len(m.schemas))
	for _, s := range m.schemas {
		schemaMap[s.field] = s
	}

	fieldMap := make(map[string]Field, len(m.Fields))
	for _, f := range m.Fields {
		fieldMap[f.Name] = f
	}

	// Check for new or changed fields
	for _, field := range m.Fields {
		schema, exists := schemaMap[field.Name]
		if !exists {
			m.AddField(&field)
			continue
		}

		filed_type, field_length := schema.parseSQLType()
		shouldChange := false

		if filed_type != field.Type.string() {
			fmt.Printf("Type mismatch on field '%s': DB = %s, Model = %s\n", field.Name, filed_type, field.Type.string())
			shouldChange = true
		}
		if !(field_length == 1 && field.Length == 0) && field_length != field.Length {
			fmt.Printf("Length mismatch on field '%s': DB = %d, Model = %d\n", field.Name, field_length, field.Length)
			shouldChange = true
		}
		if schema.defaultVal.String != field.DefaultValue {
			fmt.Printf("Default value mismatch on field '%s': DB = '%s', Model = '%s'\n", field.Name, schema.defaultVal.String, field.DefaultValue)
			shouldChange = true
		}
		if schema.nullable == "YES" && !field.Nullable {
			fmt.Printf("Nullable mismatch on field '%s': DB = YES, Model = NOT NULL\n", field.Name)
			shouldChange = true
		}
		if schema.nullable == "NO" && field.Nullable {
			fmt.Printf("Nullable mismatch on field '%s': DB = NO, Model = NULL\n", field.Name)
			shouldChange = true
		}
		if schema.extra == "auto_increment" && !field.AutoIncrement {
			fmt.Printf("Auto-increment mismatch on field '%s': DB = auto_increment, Model = not auto_increment\n", field.Name)
			shouldChange = true
		}
		switch schema.key {
		case "PRI":
			if !field.Index.PrimaryKey {
				fmt.Printf("Index mismatch on field '%s': DB = Primary Key, Model = Primary Key Removed\n", field.Name)
				shouldChange = true
			}
		case "UNI":
			if !field.Index.Unique {
				fmt.Printf("Index mismatch on field '%s': DB = Unique, Model = Unique Removed\n", field.Name)
				shouldChange = true
			}
		case "MUL":
			if !field.Index.Index {
				fmt.Printf("Index mismatch on field '%s': DB = Indexed (MUL), Model = Index Removed\n", field.Name)
				shouldChange = true
			}
		default:
			if field.Index.PrimaryKey {
				fmt.Printf("Index mismatch on field '%s': DB = None, Model = Primary Key Added\n", field.Name)
				shouldChange = true
			} else if field.Index.Unique {
				fmt.Printf("Index mismatch on field '%s': DB = None, Model = Unique Added\n", field.Name)
				shouldChange = true
			} else if field.Index.Index {
				fmt.Printf("Index mismatch on field '%s': DB = None, Model = Index Added\n", field.Name)
				shouldChange = true
			}
		}

		if shouldChange {
			m.UpdateField(&field)
		}
	}

	// Check for fields to delete
	for _, schema := range m.schemas {
		if _, exists := fieldMap[schema.field]; !exists {
			fmt.Printf("Do you want to delete %s (y/n): ", schema.field)
			reader := bufio.NewReader(os.Stdin)
			if input, err := reader.ReadString('\n'); err == nil && strings.TrimSpace(input) == "y" {
				m.DropField(schema.field)
			} else if err != nil {
				fmt.Printf("Error Getting Input: %s", err.Error())
			} else {
				fmt.Printf("Skipping the Deleting the field: %s\n", schema.field)
			}
		}
	}
}

func (m *Struct) CreateTableIfNotExists() {
	if len(m.schemas) > 0 { // if the lenth is more that 0 that means talbe is already created and no need to create it again instead we should focus on updating it
		fmt.Println(m.TableName, " is already created. So heading towards to check for the table update instead")
		m.SyncTableSchema()
		return
	}
	sql := "CREATE TABLE IF NOT EXISTS " + m.TableName + " (\n"
	fieldDefs := []string{}

	for _, field := range m.Fields {
		fieldDefs = append(fieldDefs, field.String())
	}

	sql += strings.Join(fieldDefs, ",\n")
	sql += "\n);"

	fmt.Println("Generated SQL:", sql) // 🔍 Always print this for debugging

	databaseObj, err := DatabaseHandler.GetDatabase()
	if err != nil {
		panic("Error getting database: " + err.Error())
	}

	_, err = databaseObj.Exec(sql)
	if err != nil {
		panic("Error creating table: " + err.Error() + "\nQuery:" + sql)
	}

	fmt.Println("ModelsHandler: Table created or already exists:", m.TableName)
}

// function to add the new field in the table
func (m *Struct) AddField(field *Field) {
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
			fmt.Println("Succesfully Altered the table ", m.TableName, " Field is: ", field.Name)
		}
	}
}

// function to change the field details
func (m *Struct) UpdateField(field *Field) {
	// ALTER TABLE `users` CHANGE `userId` `userId` INT(30) NOT NULL AUTO_INCREMENT;
	response := "ALTER TABLE `" + m.TableName + "` "
	response += "CHANGE `" + field.Name + "` " + field.columnDefinition() + ";"

	if databaseObj, err := DatabaseHandler.GetDatabase(); err != nil {
		panic("Error While Changing Field" + err.Error())
	} else {
		if _, sql_err := databaseObj.Exec(response); sql_err != nil {
			panic("Error While Changing the Table Field" + sql_err.Error() + "SQL QUERY: " + response)
		} else {
			fmt.Println("Succesfully Changed the table ", m.TableName, " Field is: ", field.Name)
		}
	}
}

func (m *Struct) DropField(fieldName string) {
	//ALTER TABLE `users` DROP `userId`;
	query := "ALTER TABLE `" + m.TableName + "` DROP `" + fieldName + "`;"
	if databaseObj, err := DatabaseHandler.GetDatabase(); err != nil {
		panic("Error While Deleting Field" + err.Error())
	} else {
		if _, sql_err := databaseObj.Exec(query); sql_err != nil {
			panic("Error While Deleting the Field" + sql_err.Error())
		} else {
			fmt.Println("Succesfully Deleted the field from the table ", m.TableName, " Field is: ", fieldName)
		}
	}
}

// here alter table will take the field attributes to alter the table details
// it will only alter each attribuite at a time
// Reserved for future: granular attribute-level alteration
func (m *Struct) AlterTable() string { return "" }

func (m *Struct) GetTableName() string {
	return m.TableName
}

// func NewField(name string, _type fieldType, lenth int, nullable bool, defaultValue string, index string) *Field {
// 	_field := Field{
// 		Name:         name,
// 		Type:         _type,
// 		Length:       lenth,
// 		Nullable:     nullable,
// 		DefaultValue: defaultValue,
// 		Index:        index,
// 	}

// 	if ModelsRegistry == nil {
// 		ModelsRegistry = make([]*Struct, 0)
// 	}
// 	return &_field
// }
