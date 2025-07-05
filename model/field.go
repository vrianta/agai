package model

/*
CREATE TABLE IF NOT EXISTS employees (
	id INT AUTO_INCREMENT,
	name VARCHAR(100),
	position VARCHAR(50),
	salary DECIMAL(10, 2),
	hire_date DATE,
	PRIMARY KEY `id` (id),
	INDEX `idx_name` (name)
);
*/

import (
	"fmt"
	"strconv"
	"time"
	"unicode"
)

// return the string of the total field expression mostly will be used for table creation
func (f *Field) String() string {
	response := f.name + " " + f.Type.string()
	if f.Length > 0 {
		response += "(" + fmt.Sprint(f.Length) + ") "
	}
	if f.Nullable {
		response += " NULL "
	} else if !f.Nullable {
		response += " NOT NULL "
	}

	if f.DefaultValue != "" {
		switch f.Type {
		case FieldTypes.String, FieldTypes.Text:
			response += "DEFAULT '" + f.DefaultValue + "' "
		case FieldTypes.Bool:
			if f.DefaultValue == "true" || f.DefaultValue == "1" {
				response += "DEFAULT TRUE "
			} else {
				response += "DEFAULT FALSE "
			}
		default:
			response += "DEFAULT " + f.DefaultValue + " "
		}
	}

	if f.Index.PrimaryKey {
		response += ", Primary Key pk_" + string(f.name) + " (" + string(f.name) + ")"
	}
	if f.Index.Index {
		response += ", INDEX idx_" + string(f.name) + " (" + string(f.name) + ")"
	}
	if f.Index.FullText {
		response += ", FULLTEXT ftxt_" + string(f.name) + " (" + string(f.name) + ")"
	}
	if f.Index.Spatial {
		response += ", SPATIAL sp_" + string(f.name) + " (" + string(f.name) + ")"
	}
	if f.Index.Unique {
		response += ", UNIQUE unq_" + string(f.name) + " (" + string(f.name) + ")"
	}

	return response
}

func (f *Field) columnDefinition() string {
	response := f.name + " " + f.Type.string()
	if f.Length > 0 {
		response += "(" + fmt.Sprint(f.Length) + ") "
	}
	if f.Nullable {
		response += " NULL "
	} else if !f.Nullable {
		response += " NOT NULL "
	}

	if f.DefaultValue != "" {
		switch f.Type {
		case FieldTypes.String, FieldTypes.Text:
			response += "DEFAULT '" + f.DefaultValue + "' "
		case FieldTypes.Bool:
			if f.DefaultValue == "true" || f.DefaultValue == "1" {
				response += "DEFAULT TRUE "
			} else {
				response += "DEFAULT FALSE "
			}
		default:
			response += "DEFAULT " + f.DefaultValue + " "
		}
	}

	return response
}

func (f *Field) addIndexStatement() string {
	var response string
	if f.Index.PrimaryKey {
		response += ", ADD Primary Key pk_" + string(f.name) + " (" + string(f.name) + ")"
	}
	if f.Index.Index {
		response += ", ADD INDEX idx_" + string(f.name) + " (" + string(f.name) + ")"
	}
	if f.Index.FullText {
		response += ", ADD FULLTEXT ftxt_" + string(f.name) + " (" + string(f.name) + ")"
	}
	if f.Index.Spatial {
		response += ", ADD SPATIAL sp_" + string(f.name) + " (" + string(f.name) + ")"
	}
	if f.Index.Unique {
		response += ", ADD UNIQUE unq_" + string(f.name) + " (" + string(f.name) + ")"
	}

	return response
}

func (f *Field) Name() string {
	return f.name
}

func (ft fieldType) IsValueCompatible(val string) bool {
	switch ft {
	case FieldTypes.Int, FieldTypes.TinyInt, FieldTypes.SmallInt, FieldTypes.MediumInt, FieldTypes.BigInt:
		_, err := strconv.Atoi(val)
		return err == nil
	case FieldTypes.Float, FieldTypes.Decimal, FieldTypes.Double, FieldTypes.Real:
		_, err := strconv.ParseFloat(val, 64)
		return err == nil
	case FieldTypes.Date, FieldTypes.Time, FieldTypes.Timestamp:
		_, err1 := time.Parse("2006-01-02", val)
		_, err2 := time.Parse("2006-01-02 15:04:05", val)
		return err1 == nil || err2 == nil
	case FieldTypes.VarChar, FieldTypes.Text, FieldTypes.String, FieldTypes.Char:
		return true
	case FieldTypes.Bool:
		return val == "0" || val == "1" || val == "true" || val == "false"
	default:
		return true
	}
}

func isAlphaNumeric(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

func (f *Field) Compare(fieldTypeStr string) bool {
	switch fieldTypeStr {
	case "TINYINT":
		switch f.Type {
		case FieldTypes.TinyInt, FieldTypes.Bool:
			return true
		}
	case "SMALLINT":
		switch f.Type {
		case FieldTypes.SmallInt:
			return true
		}
	case "MEDIUMINT":
		switch f.Type {
		case FieldTypes.MediumInt:
			return true
		}
	case "INT", "INTEGER":
		switch f.Type {
		case FieldTypes.Int:
			return true
		}
	case "BIGINT":
		switch f.Type {
		case FieldTypes.BigInt:
			return true
		}
	case "VARCHAR":
		switch f.Type {
		case FieldTypes.VarChar, FieldTypes.String:
			return true
		}
	case "CHAR":
		switch f.Type {
		case FieldTypes.Char, FieldTypes.String:
			return true
		}
	case "TEXT":
		switch f.Type {
		case FieldTypes.Text, FieldTypes.String:
			return true
		}
	case "LONGTEXT":
		switch f.Type {
		case FieldTypes.LongText, FieldTypes.JSON:
			return true
		}
	case "BOOL", "BOOLEAN":
		switch f.Type {
		case FieldTypes.Bool, FieldTypes.TinyInt:
			return true
		}
	case "DECIMAL":
		switch f.Type {
		case FieldTypes.Decimal:
			return true
		}
	case "FLOAT":
		switch f.Type {
		case FieldTypes.Float:
			return true
		}
	case "DOUBLE", "REAL":
		switch f.Type {
		case FieldTypes.Double, FieldTypes.Real:
			return true
		}
	case "JSON":
		switch f.Type {
		case FieldTypes.JSON:
			return true
		}
	case "BLOB":
		switch f.Type {
		case FieldTypes.Blob:
			return true
		}
	case "DATE":
		switch f.Type {
		case FieldTypes.Date:
			return true
		}
	case "TIME":
		switch f.Type {
		case FieldTypes.Time:
			return true
		}
	case "TIMESTAMP":
		switch f.Type {
		case FieldTypes.Timestamp:
			return true
		}
	case "ENUM":
		switch f.Type {
		case FieldTypes.Enum:
			return true
		}
	case "UUID":
		switch f.Type {
		case FieldTypes.UUID:
			return true
		}
	default:
		return false
	}

	return false
}
