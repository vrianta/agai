package structparser

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var args map[string]string = make(map[string]string) // args
var args_scanned bool = false                        // it keeps a track if we already gone through all the args

// parseArg processes CLI arguments based on the "flag" tag.
func parseArg(flagKey string, fieldVal *reflect.Value, fieldType *reflect.StructField) error {
	argValue, exists := lookupArg(flagKey)
	if !exists || argValue == "" {
		return fmt.Errorf("Did not find the element %s in the argument", flagKey) // Skip if the CLI argument wasn't provided
	}

	if !fieldVal.CanSet() {
		return fmt.Errorf("field %s is not public", fieldType.Name)
	}

	switch fieldVal.Kind() {
	case reflect.Pointer:
		// 1. Allocate memory if the pointer is currently nil
		if fieldVal.IsNil() {
			fieldVal.Set(reflect.New(fieldVal.Type().Elem()))
		}

		// 2. Dereference the pointer to get the underlying value
		elem := fieldVal.Elem()

		// 3. Parse into the underlying value recursively
		return parseArg(flagKey, &elem, fieldType)
	case reflect.String:
		fieldVal.SetString(argValue)

	case reflect.Bool:
		b, err := strconv.ParseBool(argValue)
		if err != nil {
			return fmt.Errorf("field %s: cannot parse arg %q as bool: %w", fieldType.Name, argValue, err)
		}
		fieldVal.SetBool(b)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if fieldType.Type.String() == "time.Duration" {
			d, err := time.ParseDuration(argValue)
			if err != nil {
				return fmt.Errorf("field %s: cannot parse arg %q as time.Duration: %w", fieldType.Name, argValue, err)
			}
			fieldVal.SetInt(int64(d))
			break
		}

		parsedInt, err := strconv.ParseInt(argValue, 10, 64)
		if err != nil {
			return fmt.Errorf("field %s: cannot parse arg %q as int: %w", fieldType.Name, argValue, err)
		}
		if fieldVal.OverflowInt(parsedInt) {
			return fmt.Errorf("field %s: arg value %s overflows %s", fieldType.Name, argValue, fieldVal.Kind())
		}
		fieldVal.SetInt(parsedInt)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		parsedUint, err := strconv.ParseUint(argValue, 10, 64)
		if err != nil {
			return fmt.Errorf("field %s: cannot parse arg %q as uint: %w", fieldType.Name, argValue, err)
		}
		if fieldVal.OverflowUint(parsedUint) {
			return fmt.Errorf("field %s: arg value %s overflows %s", fieldType.Name, argValue, fieldVal.Kind())
		}
		fieldVal.SetUint(parsedUint)

	case reflect.Float32, reflect.Float64:
		parsedFloat, err := strconv.ParseFloat(argValue, 64)
		if err != nil {
			return fmt.Errorf("field %s: cannot parse arg %q as float: %w", fieldType.Name, argValue, err)
		}
		if fieldVal.OverflowFloat(parsedFloat) {
			return fmt.Errorf("field %s: arg value %s overflows %s", fieldType.Name, argValue, fieldVal.Kind())
		}
		fieldVal.SetFloat(parsedFloat)

	default:
		return fmt.Errorf("field %s has unsupported type %s for arg values", fieldType.Name, fieldVal.Kind())
	}

	return nil
}

// ParseOSArgs scans os.Args and caches command-line arguments for Parse and New.
//
// It supports --flag=value, -flag=value, --flag value, -flag value, and
// boolean flags such as --debug. Arguments without a leading dash are ignored.
// Calling ParseOSArgs more than once has no effect after the first scan.
func ParseOSArgs() {
	if args_scanned {
		return
	}

	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]

		// Ignore arguments that don't start with "-" or "--"
		if !strings.HasPrefix(arg, "-") {
			continue
		}

		// Trim leading dashes (handles both "-" and "--")
		trimmed := strings.TrimLeft(arg, "-")

		// Case 1: Key=Value format (e.g., --port=8080 or -port=8080)
		if key, value, found := strings.Cut(trimmed, "="); found {
			args[key] = value
			continue
		}

		// Case 2: Space-separated format (e.g., --port 8080)
		if i+1 < len(os.Args) && !strings.HasPrefix(os.Args[i+1], "-") {
			args[trimmed] = os.Args[i+1]
			i++ // Skip next index since it was consumed as a value
			continue
		}

		// Case 3: Boolean flag format (e.g., --debug or -v)
		args[trimmed] = "true"
	}

	args_scanned = true
}

// lookupArg now becomes a lightweight O(1) query function
func lookupArg(flagName string) (string, bool) {
	if !args_scanned {
		ParseOSArgs()
	}

	flagV, ok := args[flagName]
	return flagV, ok
}
