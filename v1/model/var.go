package model

const (
	String fieldType = iota
	Text
	VarChar
	Int
	Float
	Decimal
	Bool
	Date
	Time
	Timestamp
	JSON
	Enum
	Binary
	UUID
	TinyInt

	// Additional SQL types
	SmallInt
	MediumInt
	BigInt
	Double
	Real
	Char
	LongText
	MediumText
	TinyText
	Blob
	TinyBlob
	MediumBlob
	LongBlob
	Set
	Year
	Geometry
	Point
	LineString
	Polygon
)

var (
	FieldTypes = struct {
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
		TinyInt   fieldType

		// New types
		SmallInt   fieldType
		MediumInt  fieldType
		BigInt     fieldType
		Double     fieldType
		Real       fieldType
		Char       fieldType
		LongText   fieldType
		MediumText fieldType
		TinyText   fieldType
		Blob       fieldType
		TinyBlob   fieldType
		MediumBlob fieldType
		LongBlob   fieldType
		Set        fieldType
		Year       fieldType
		Geometry   fieldType
		Point      fieldType
		LineString fieldType
		Polygon    fieldType
	}{
		String:    String,
		Text:      Text,
		VarChar:   VarChar,
		Int:       Int,
		Float:     Float,
		Decimal:   Decimal,
		Bool:      Bool,
		Date:      Date,
		Time:      Time,
		Timestamp: Timestamp,
		JSON:      JSON,
		Enum:      Enum,
		Binary:    Binary,
		UUID:      UUID,
		TinyInt:   TinyInt,

		SmallInt:   SmallInt,
		MediumInt:  MediumInt,
		BigInt:     BigInt,
		Double:     Double,
		Real:       Real,
		Char:       Char,
		LongText:   LongText,
		MediumText: MediumText,
		TinyText:   TinyText,
		Blob:       Blob,
		TinyBlob:   TinyBlob,
		MediumBlob: MediumBlob,
		LongBlob:   LongBlob,
		Set:        Set,
		Year:       Year,
		Geometry:   Geometry,
		Point:      Point,
		LineString: LineString,
		Polygon:    Polygon,
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

	ModelsRegistry = map[string]*meta{}

	sqlKeywords = map[string]bool{
		"ADD": true, "ALL": true, "ALTER": true, "AND": true, "ANY": true,
		"AS": true, "ASC": true, "BACKUP": true, "BETWEEN": true, "CASE": true,
		"CHECK": true, "COLUMN": true, "CONSTRAINT": true, "CREATE": true,
		"DATABASE": true, "DEFAULT": true, "DELETE": true, "DESC": true,
		"DISTINCT": true, "DROP": true, "EXEC": true, "EXISTS": true,
		"FOREIGN": true, "FROM": true, "FULL": true, "GROUP": true, "HAVING": true,
		"IN": true, "INDEX": true, "INNER": true, "INSERT": true, "INTO": true,
		"IS": true, "JOIN": true, "KEY": true, "LEFT": true, "LIKE": true,
		"LIMIT": true, "NOT": true, "NULL": true, "OR": true, "ORDER": true,
		"OUTER": true, "PRIMARY": true, "PROCEDURE": true, "RIGHT": true,
		"ROWNUM": true, "SELECT": true, "SET": true, "TABLE": true, "TOP": true,
		"TRUNCATE": true, "UNION": true, "UNIQUE": true, "UPDATE": true,
		"VALUES": true, "VIEW": true, "WHERE": true, "WITH": true,
	}

	//
	// jsonStoreMu      sync.RWMutex
	componentsDir string = "./components"
	// warnedMissingDir = false
	// wb               = sync.WaitGroup{}

	initialsed = false
)
