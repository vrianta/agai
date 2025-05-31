package Template

import (
	"html/template"
	"time"
)

type (
	Response map[string]any

	Struct struct {
		Uri          string            // path of the template file
		LastModified time.Time         // date when the file last modified
		Data         template.Template // template data of the file before modified
	}
)
