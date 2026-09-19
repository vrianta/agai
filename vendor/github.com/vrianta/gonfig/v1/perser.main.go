package structparser

import (
	"errors"
	"fmt"
	"reflect"
)

// New creates and populates a configuration value of type T.
//
// The value is populated from environment variables, command-line arguments,
// and default tags in that order. If crashOnFail is true, a missing required
// field causes New to panic. When crashOnFail is false, parsing errors are not
// returned; use Parse when the caller needs to inspect errors.
func New[T any](crashOnFail bool) T {
	var cfg T
	_, _ = Parse(&cfg, crashOnFail)
	return cfg
}

// Parse populates the exported fields of a struct from its configuration tags.
//
// ctx must be a non-nil pointer to a struct. Environment variables take
// precedence over command-line arguments, which take precedence over default
// values. Nested structs are processed recursively, and pointer fields are
// allocated when a value is provided.
//
// Parse returns an error for invalid input, unexported fields, and missing
// required fields. If crashOnFail is true, a missing required field causes a
// panic instead of returning an error.
func Parse[T any](ctx *T, crashOnFail bool) (*T, error) {

	val := reflect.ValueOf(ctx)

	// Ensure we were passed a non-nil pointer
	if val.Kind() != reflect.Pointer || val.IsNil() {
		return nil, errors.New("must pass a non-nil pointer to a struct")
	}

	elem := val.Elem()
	if elem.Kind() != reflect.Struct {
		return nil, errors.New("provided value is not a pointer to a struct")
	}

	_, printHelp := lookupArg("help")
	var records []help_record

	if err := parseStructValue(elem, crashOnFail, printHelp, &records); err != nil {
		return nil, err
	}

	if printHelp {
		help(records)
	}

	return ctx, nil
}

func parseStructValue(val reflect.Value, crashOnFail, printHelp bool, records *[]help_record) error {
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		fieldVal := val.Field(i)
		fieldType := typ.Field(i)

		// Handle nested struct recursively (must pass addressable value)
		if fieldVal.Kind() == reflect.Struct {
			if fieldVal.CanAddr() {
				if err := parseStructValue(fieldVal, crashOnFail, printHelp, records); err != nil {
					return err
				}
			}
			continue
		}

		// Ensure unexported/private fields are rejected early
		if !fieldVal.CanSet() {
			return fmt.Errorf("field %s is not public", fieldType.Name)
		}

		rec := help_record{field: fieldType.Name}
		var matched bool

		// 1. Process Environment Variables
		if envKey, ok := fieldType.Tag.Lookup("env"); ok && envKey != "" {
			rec.env = envKey
			if err := parseEnv(envKey, &fieldVal, &fieldType); err == nil {
				matched = true
			}
		}

		// 2. Process Command-Line Arguments
		if !matched {
			if argName, ok := fieldType.Tag.Lookup("arg"); ok && argName != "" {
				rec.args = argName
				if err := parseArg(argName, &fieldVal, &fieldType); err == nil {
					matched = true
				}
			}
		}

		// 3. Fallback to Default Value
		if !matched {
			if defVal, ok := fieldType.Tag.Lookup("default"); ok && defVal != "" {
				rec.def = defVal
				_ = parseDefault(defVal, &fieldVal, &fieldType)
			}
		}

		// Extract description tag for help text
		if desc, ok := fieldType.Tag.Lookup("description"); ok {
			rec.description = desc
		}

		// 4. Validate Required Constraints
		if _, isRequired := fieldType.Tag.Lookup("required"); isRequired {
			if fieldVal.IsZero() {
				errMsg := fmt.Sprintf("field %s is required but has no value", fieldType.Name)
				if crashOnFail && !printHelp {
					panic(errMsg)
				}
				return errors.New(errMsg)
			}
		}

		if printHelp {
			*records = append(*records, rec)
		}
	}

	return nil
}
