package model

// For the sake of go comunity who were so pissed off because of the capital lattern, becuse I think they are allergic to it,
// I Decided to give this file name with smaller latter and _ becuase I do not want to make the kids more raged on this, and wasting my time
// I hope this will make them happy and they will stop crying about it

func (f fieldType) string() string {
	switch f {
	case FieldsTypes.String, FieldsTypes.VarChar:
		return "VARCHAR"
	case FieldsTypes.Text:
		return "TEXT"
	case FieldsTypes.Int:
		return "INT"
	case FieldsTypes.Float:
		return "FLOAT"
	case FieldsTypes.Decimal:
		return "DECIMAL(10,2)"
	case FieldsTypes.Bool:
		return "TINYINT"
	case FieldsTypes.Date:
		return "DATE"
	case FieldsTypes.Time:
		return "DATETIME"
	case FieldsTypes.Timestamp:
		return "TIMESTAMP"
	case FieldsTypes.JSON:
		return "LONGTEXT"
	case FieldsTypes.Enum:
		return "ENUM" // You can customize enum values at the field level
	case FieldsTypes.Binary:
		return "BLOB"
	case FieldsTypes.UUID:
		return "CHAR(36)" // UUIDs typically stored as 36-char strings
	default:
		return "TEXT" // Safe fallback
	}
}
