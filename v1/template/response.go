package template

import (
	"fmt"
	"html/template"
	"reflect"
	"strings"
)

var ReponseFuncMaps = template.FuncMap{
	"upper": strings.ToUpper,
	"lower": strings.ToLower,
	"strlen": func(s string) int {
		return len(s)
	},
	"len": func(v any) int {
		rv := reflect.ValueOf(v)
		switch rv.Kind() {
		case reflect.Array, reflect.Slice, reflect.Map, reflect.Chan, reflect.String:
			return rv.Len()
		default:
			return 0
		}
	},
	"print": func(data any) string {
		return fmt.Sprintln(data)
	},
}
