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

import "fmt"

// returnt he string of the total field expression mostly will be used for table creation
func (f *Field) String() string {
	response := f.Name + " " + f.Type.string()
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
		case FieldsTypes.String, FieldsTypes.Text:
			response += "DEFAULT '" + f.DefaultValue + "' "
		case FieldsTypes.Bool:
			if f.DefaultValue == "true" || f.DefaultValue == "1" {
				response += "DEFAULT TRUE "
			} else {
				response += "DEFAULT FALSE "
			}
		default:
			// For numeric types, we assume the default value is a number
			response += "DEFAULT " + f.DefaultValue + " "
		}
	}

	if f.Index.PrimaryKey { // it is a primary key
		response += ", Primary Key pk_" + string(f.Name) + " (" + string(f.Name) + ")"
	}
	if f.Index.Index {
		response += ", INDEX idx_" + string(f.Name) + " (" + string(f.Name) + ")"
	}
	if f.Index.FullText {
		response += ", FULLTEXT ftxt_" + string(f.Name) + " (" + string(f.Name) + ")"
	}
	if f.Index.Spatial {
		response += ", SPATIAL sp_" + string(f.Name) + " (" + string(f.Name) + ")"
	}
	if f.Index.Unique {
		response += ", UNIQUE unq_" + string(f.Name) + " (" + string(f.Name) + ")"
	}

	return response
}

func (f *Field) columnDefinition() string {
	response := f.Name + " " + f.Type.string()
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
		case FieldsTypes.String, FieldsTypes.Text:
			response += "DEFAULT '" + f.DefaultValue + "' "
		case FieldsTypes.Bool:
			if f.DefaultValue == "true" || f.DefaultValue == "1" {
				response += "DEFAULT TRUE "
			} else {
				response += "DEFAULT FALSE "
			}
		default:
			// For numeric types, we assume the default value is a number
			response += "DEFAULT " + f.DefaultValue + " "
		}
	}

	return response
}

func (f *Field) addIndexStatement() string {
	// ADD INDEX (`newel`);
	var response string
	if f.Index.PrimaryKey { // it is a primary key
		response += ", ADD Primary Key pk_" + string(f.Name) + " (" + string(f.Name) + ")"
	}
	if f.Index.Index {
		response += ", ADD INDEX idx_" + string(f.Name) + " (" + string(f.Name) + ")"
	}
	if f.Index.FullText {
		response += ", ADD FULLTEXT ftxt_" + string(f.Name) + " (" + string(f.Name) + ")"
	}
	if f.Index.Spatial {
		response += ", ADD SPATIAL sp_" + string(f.Name) + " (" + string(f.Name) + ")"
	}
	if f.Index.Unique {
		response += ", ADD UNIQUE unq_" + string(f.Name) + " (" + string(f.Name) + ")"
	}

	return response
}
