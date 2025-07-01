package models_handler

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
		model *Struct

		// WHERE clause
		whereClauses []string
		whereArgs    []any
		lastColumn   string

		// SET clause for update
		setClauses []string
		setArgs    []any
		lastSet    string
		groupBy    string

		// Other options
		limit   int
		offset  int
		orderBy string

		operation string // "select", "delete", "update"
	}
)
