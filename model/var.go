package model

var (
	FieldsTypes = struct {
		String    fieldType
		Text      fieldType
		VarChar   fieldType
		Int       fieldType
		Float     fieldType
		Decimal   fieldType
		Bool      fieldType
		Date      fieldType
		Time      fieldType
		Timestamp fieldType
		JSON      fieldType
		Enum      fieldType
		Binary    fieldType
		UUID      fieldType
	}{
		String:    0,
		Text:      1,
		VarChar:   2,
		Int:       3,
		Float:     4,
		Decimal:   5,
		Bool:      6,
		Date:      7,
		Time:      8,
		Timestamp: 9,
		JSON:      10,
		Enum:      11,
		Binary:    12,
		UUID:      13,
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

	ModelsRegistry = map[string]*Struct{}
)
