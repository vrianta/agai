package Models

import "database/sql"

type (
	fieldType uint16
	Fields    map[string]fieldType

	Struct struct {
		TableName string  // Name of the table in the database
		Fields    []Field // Map of field names to their types
		schemas   []schema
	}

	Index struct {
		PrimaryKey bool
		Unique     bool
		Index      bool
		FullText   bool
		Spatial    bool
	}

	Field struct {
		Name          string
		Type          fieldType
		Length        int
		Nullable      bool
		DefaultValue  string
		AutoIncrement bool
		Index         Index // Index type (e.g., "UNIQUE", "INDEX")
	}

	schema struct { // Scema is to hold the table scema which is available in the database
		field, fieldType, nullable, key, extra string
		defaultVal                             sql.NullString
	}
)
