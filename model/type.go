package model

import "database/sql"

type (
	fieldType    uint16
	FieldTypes   map[string]fieldType
	FieldTypeset map[string]*Field
	Result       map[string]any
	Results      map[any]Result

	// indexInfo struct {
	// 	ColumnName string
	// 	IndexName  string
	// 	NonUnique  bool
	// }

	Field struct {
		Name          string
		Type          fieldType
		Length        int
		Nullable      bool
		DefaultValue  string
		AutoIncrement bool
		Index         Index // Index type (e.g., "UNIQUE", "INDEX")
	}

	Table struct {
		TableName   string       // Name of the table in the database
		FieldTypes  FieldTypeset // Map of field names to their types
		schemas     []schema
		Initialised bool   // Flag to check if the model is initialised
		primary     *Field // name of the primary elemet
		// indexes     map[string]indexInfo // columnName -> index info
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

	queryBuilder struct {
		model *Table

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

		operation           string // "select", "delete", "update"
		InsertRowFieldTypes map[string]any
	}

	// InsertRowBuilder is a dedicated struct for InsertRow operations (CREATE), separate from the general queryBuilder struct.
	InsertRowBuilder struct {
		model               *Table
		InsertRowFieldTypes map[string]any
		lastSet             string
	}
)
