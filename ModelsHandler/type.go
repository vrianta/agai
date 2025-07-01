package models

import "database/sql"

type (
	fieldType uint16
	Fields    map[string]fieldType

	Struct struct {
		TableName string           // Name of the table in the database
		fields    map[string]Field // Map of field names to their types
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
		value         any   // a variable where the value of the field will be stored
	}

	schema struct { // Scema is to hold the table scema which is available in the database
		field, fieldType, nullable, key, extra string
		defaultVal                             sql.NullString
	}

	Query struct {
		model      *Struct
		conditions []string
		args       []any
		limit      int
		offset     int
		lastColumn string
	}
)
