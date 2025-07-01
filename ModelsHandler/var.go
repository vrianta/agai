package models

var (
	FieldsTypes = struct {
		String  fieldType
		Text    fieldType
		VarChar fieldType
		Int     fieldType
		Float   fieldType
		Bool    fieldType
		Date    fieldType
		Time    fieldType
		JSON    fieldType
		Decimal fieldType
	}{
		String:  0, // String type field
		Text:    1, // Text type field
		VarChar: 2, // Variable character field
		Int:     3, // Integer type field
		Float:   4, // Float type field
		Bool:    5, // Boolean type field
		Date:    6, // Date type field
		Time:    7, // Time type field
		JSON:    8, // JSON type field
		Decimal: 9,
	}

	// Indexes = struct {
	// 	PrimaryKey Index
	// 	Unique     Index
	// 	Index      Index
	// 	FullText   Index
	// 	Spatial    Index
	// }{
	// 	PrimaryKey: "Primary Key",
	// 	Unique:     "UNIQUE",   // Unique index
	// 	Index:      "INDEX",    // Regular index
	// 	FullText:   "FULLTEXT", // Full-text index
	// 	Spatial:    "SPATIAL",  // Spatial index
	// }

	ModelsRegistry = []*Struct{}
)
