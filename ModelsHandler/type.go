package Models

type (
	fieldType uint16

	Struct struct {
		TableName  string // Name of the table in the database
		PrimaryKey string // Primary key for the model

		Fields map[string]fieldType // Map of field names to their types
	}
)
