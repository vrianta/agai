package Router

import (
	"sync"
	"time"

	"github.com/vrianta/Server/Controller"
)

type (
	Type map[string]*Controller.Struct // Type for routes, mapping URL paths to Controller structs

	Struct struct {
		sessions sync.Map
		routes   Type
	}

	FileInfo struct {
		Uri          string    // path of the template file
		LastModified time.Time // date when the file last modified
		Data         string    // template data of the file before modified
	}
)
