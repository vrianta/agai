package model

// For the sake of go comunity who were so pissed off because of the capital lattern, becuse I think they are allergic to it,
// I Decided to give this file name with smaller latter and _ becuase I do not want to make the kids more raged on this, and wasting my time
// I hope this will make them happy and they will stop crying about it

func (f fieldType) string() string {
	switch f {
	case FieldTypes.String, FieldTypes.VarChar:
		return "VARCHAR"
	case FieldTypes.Text:
		return "TEXT"
	case FieldTypes.Int:
		return "INT"
	case FieldTypes.Float:
		return "FLOAT"
	case FieldTypes.Decimal:
		return "DECIMAL(10,2)"
	case FieldTypes.Bool:
		return "BOOLEAN"
	case FieldTypes.TinyInt:
		return "TINYINT"
	case FieldTypes.Date:
		return "DATE"
	case FieldTypes.Time:
		return "DATETIME"
	case FieldTypes.Timestamp:
		return "TIMESTAMP"
	case FieldTypes.JSON:
		return "LONGTEXT"
	case FieldTypes.Enum:
		return "ENUM" // You can customize enum values at the field level
	case FieldTypes.Binary:
		return "BLOB"
	case FieldTypes.UUID:
		return "CHAR(36)" // UUIDs typically stored as 36-char strings
	default:
		return "TEXT" // Safe fallback
	}
}

func (ft fieldType) IsNumeric() bool {
	switch ft {
	case FieldTypes.Bool,
		FieldTypes.Int,
		FieldTypes.TinyInt,
		FieldTypes.Float,
		FieldTypes.Decimal:
		return true
	default:
		return false
	}
}
