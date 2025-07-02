package router

import "sync"

var (
	fileCache  sync.Map           // map[string]FileInfo
	routeTable = make(routes, 50) // map[string]*Controller.Struct
)
