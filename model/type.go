package model

import "database/sql"

type (
	fieldType uint16
	Fields    map[string]fieldType
	FieldMap  map[string]*Field
	Result    map[string]any
	Results   map[any]Result

	indexInfo struct {
		ColumnName string
		IndexName  string
		NonUnique  bool
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

	Struct struct {
		TableName   string   // Name of the table in the database
		fields      FieldMap // Map of field names to their types
		schemas     []schema
		Initialised bool                 // Flag to check if the model is initialised
		primary     *Field               // name of the primary elemet
		indexes     map[string]indexInfo // columnName -> index info
	}

	Index struct {
		PrimaryKey bool
		Unique     bool
		Index      bool
		FullText   bool
		Spatial    bool
	}

	schema struct {
		field      string
		fieldType  string
		nullable   string
		key        string
		extra      string
		defaultVal sql.NullString

		// Add these for precise index detection (from `information_schema.statistics`)
		// indexName string
		isunique  bool
		isindex   bool
		isprimary bool
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

		operation    string // "select", "delete", "update"
		insertFields map[string]any
	}

	// InsertQuery is a dedicated struct for insert operations (CREATE), separate from the general Query struct.
	InsertQuery struct {
		model        *Struct
		insertFields map[string]any
		lastSet      string
	}
)
