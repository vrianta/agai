package models

// For the sake of go comunity who were so pissed off because of the capital lattern, becuse I think they are allergic to it,
// I Decided to give this file name with smaller latter and _ becuase I do not want to make the kids more raged on this, and wasting my time
// I hope this will make them happy and they will stop crying about it

func (f fieldType) string() string {
	switch f {
	case FieldsTypes.String:
		return "CHAR"
	case FieldsTypes.Text:
		return "TEXT"
	case FieldsTypes.VarChar:
		return "VARCHAR"
	case FieldsTypes.Int:
		return "INT"
	case FieldsTypes.Float:
		return "FLOAT"
	case FieldsTypes.Bool:
		return "BOOL"
	case FieldsTypes.Date:
		return "DATE"
	case FieldsTypes.Time:
		return "TIME"
	case FieldsTypes.JSON:
		return "JSON"
	case FieldsTypes.Decimal:
		return "DECIMAL"
	default:
		return "UNKNOWN"
	}
}
