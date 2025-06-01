package Router

import (
	"sync"
	"time"

	"github.com/vrianta/Server/Controller"
	"github.com/vrianta/Server/Session"
)

type (
	Type map[string]*Controller.Struct // Type for routes, mapping URL paths to Controller structs

	Struct struct {
		sessions     map[string]*Session.Struct // Map to hold user sessions, key is session ID
		sessionMutex sync.RWMutex               // Mutex for thread-safe session access
		routes       Type
	}

	FileInfo struct {
		Uri          string    // path of the template file
		LastModified time.Time // date when the file last modified
		Data         string    // template data of the file before modified
	}
)
