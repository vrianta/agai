package Models

var (
	FieldsTypes = struct {
		String  fieldType
		VarChar fieldType
		Int     fieldType
		Float   fieldType
		Bool    fieldType
		Date    fieldType
		Time    fieldType
		JSON    fieldType
	}{
		String:  0, // String type field
		VarChar: 1, // Variable character field
		Int:     2, // Integer type field
		Float:   3, // Float type field
		Bool:    4, // Boolean type field
		Date:    5, // Date type field
		Time:    6, // Time type field
		JSON:    7, // JSON type field

	}
)
