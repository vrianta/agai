package RenderEngine

import (
	"bytes"
	"sync"

	"github.com/vrianta/Server/Template"
)

var (
	templateRecords = make(map[string]Template.Struct) // keep the reocrd of all the templated which are already templated
	// Use a sync.Pool for bytes.Buffer to reduce allocations
	bufPool = sync.Pool{
		New: func() interface{} { return new(bytes.Buffer) },
	}
)
