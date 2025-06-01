package RenderEngine

import (
	"bytes"
	"sync"

	"github.com/vrianta/Server/Template"
)

var (
	templateRecords      = make(map[string]Template.Struct) // keep the record of all the templates which are already templated
	templateRecordsMutex = &sync.RWMutex{}
	// Use a sync.Pool for bytes.Buffer to reduce allocations
	bufPool = sync.Pool{
		New: func() interface{} { return new(bytes.Buffer) },
	}
)
