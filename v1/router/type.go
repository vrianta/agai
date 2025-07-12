package router

import (
	"time"

	Controller "github.com/vrianta/agai/v1/controller"
)

type (
	route struct {
		path             string
		controllerObject Controller.Context
	}
	routes map[string]*Controller.Context // Type for routes, mapping URL paths to Controller structs

	Struct struct {
		defaultRoute string
	}

	FileCacheEntry struct {
		Uri          string    // path of the template file
		LastModified time.Time // date when the file last modified
		Data         string    // template data of the file before modified
	}
)
