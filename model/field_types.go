package model

// For the sake of go comunity who were so pissed off because of the capital lattern, becuse I think they are allergic to it,
// I Decided to give this file name with smaller latter and _ becuase I do not want to make the kids more raged on this, and wasting my time
// I hope this will make them happy and they will stop crying about it

func (f fieldType) string() string {
	switch f {
	case FieldTypesTypes.String, FieldTypesTypes.VarChar:
		return "VARCHAR"
	case FieldTypesTypes.Text:
		return "TEXT"
	case FieldTypesTypes.Int:
		return "INT"
	case FieldTypesTypes.Float:
		return "FLOAT"
	case FieldTypesTypes.Decimal:
		return "DECIMAL(10,2)"
	case FieldTypesTypes.Bool:
		return "TINYINT"
	case FieldTypesTypes.Date:
		return "DATE"
	case FieldTypesTypes.Time:
		return "DATETIME"
	case FieldTypesTypes.Timestamp:
		return "TIMESTAMP"
	case FieldTypesTypes.JSON:
		return "LONGTEXT"
	case FieldTypesTypes.Enum:
		return "ENUM" // You can customize enum values at the field level
	case FieldTypesTypes.Binary:
		return "BLOB"
	case FieldTypesTypes.UUID:
		return "CHAR(36)" // UUIDs typically stored as 36-char strings
	default:
		return "TEXT" // Safe fallback
	}
}
