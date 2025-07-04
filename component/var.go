package component

import "sync"

var (
	jsonStore        storage // store all the tables
	jsonStoreMu      sync.RWMutex
	componentsDir    = "./components"
	warnedMissingDir = false
)
