package Router

import "sync"

var (
	fileInfo sync.Map // map[string]FileInfo
	routes   Routes   // map[string]*Controller.Struct
)
