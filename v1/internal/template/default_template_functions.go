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
	"include": func(template_idx string, c any) template.HTML {
		if t, ok := templateComponents[template_idx]; !ok {
			return template.HTML("No Template Found: " + template_idx)
		} else {
			if data, err := t.Execute(c); err != nil {
				return template.HTML("Failed to Execute the template: " + template_idx + " Error: " + err.Error())
			} else {
				// fmt.Println(string(data))
				return template.HTML(string(data))
			}
		}
	},
}
