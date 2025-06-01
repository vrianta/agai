package Session

import "sync"

var (
	all   = make(map[string]*Struct) // Map to hold user sessions, key is session ID
	mutex = sync.RWMutex{}           // Mutex for thread-safe session access
)
