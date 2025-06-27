package Models

type (
	fieldType uint16

	Struct struct {
		PrimaryKey string               // Primary key for the model
		TableName  string               // Name of the table in the database
		Fields     map[string]fieldType // Map of field names to their types
	}
)
