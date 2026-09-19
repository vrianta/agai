package structparser

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

// parseDefault parses only the Default tag which is default
func parseDefault(defaultValue string, fieldVal *reflect.Value, fieldType *reflect.StructField) error {

	switch fieldVal.Kind() {
	case reflect.Pointer:
		// 1. Allocate memory if the pointer is currently nil
		if fieldVal.IsNil() {
			fieldVal.Set(reflect.New(fieldVal.Type().Elem()))
		}

		// 2. Dereference the pointer to get the underlying value
		elem := fieldVal.Elem()

		// 3. Parse into the underlying value recursively
		return parseDefault(defaultValue, &elem, fieldType)

	case reflect.String:
		fieldVal.SetString(defaultValue)

	case reflect.Bool:
		b, err := strconv.ParseBool(defaultValue)
		if err != nil {
			return fmt.Errorf("field %s: cannot parse default %q as bool: %w", fieldType.Name, defaultValue, err)
		}
		fieldVal.SetBool(b)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		// Check for time.Duration explicitly (type alias for int64)
		if fieldVal.Type() == reflect.TypeOf(time.Duration(0)) {
			d, err := time.ParseDuration(defaultValue)
			if err != nil {
				return fmt.Errorf("field %s: cannot parse default %q as time.Duration: %w", fieldType.Name, defaultValue, err)
			}
			fieldVal.SetInt(int64(d))
			break
		}

		parsedInt, err := strconv.ParseInt(defaultValue, 10, 64)
		if err != nil {
			return fmt.Errorf("field %s: cannot parse default %q as int: %w", fieldType.Name, defaultValue, err)
		}
		if fieldVal.OverflowInt(parsedInt) {
			return fmt.Errorf("field %s: default value %s overflows %s", fieldType.Name, defaultValue, fieldVal.Kind())
		}
		fieldVal.SetInt(parsedInt)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		parsedUint, err := strconv.ParseUint(defaultValue, 10, 64)
		if err != nil {
			return fmt.Errorf("field %s: cannot parse default %q as uint: %w", fieldType.Name, defaultValue, err)
		}
		if fieldVal.OverflowUint(parsedUint) {
			return fmt.Errorf("field %s: default value %s overflows %s", fieldType.Name, defaultValue, fieldVal.Kind())
		}
		fieldVal.SetUint(parsedUint)

	case reflect.Float32, reflect.Float64:
		parsedFloat, err := strconv.ParseFloat(defaultValue, 64)
		if err != nil {
			return fmt.Errorf("field %s: cannot parse default %q as float: %w", fieldType.Name, defaultValue, err)
		}
		if fieldVal.OverflowFloat(parsedFloat) {
			return fmt.Errorf("field %s: default value %s overflows %s", fieldType.Name, defaultValue, fieldVal.Kind())
		}
		fieldVal.SetFloat(parsedFloat)

	default:
		return fmt.Errorf("field %s has unsupported type %s for default values", fieldType.Name, fieldVal.Kind())
	}

	return nil
}
