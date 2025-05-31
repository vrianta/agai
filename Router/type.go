package Router

import (
	"time"

	"github.com/vrianta/Server/Controller"
	Session "github.com/vrianta/Server/Session"
)

type (
	Type map[string]Controller.Struct

	Struct struct {
		sessions map[string]Session.Struct
		routes   Type
	}

	FileInfo struct {
		Uri          string    // path of the template file
		LastModified time.Time // date when the file last modified
		Data         string    // template data of the file before modified
	}
)
