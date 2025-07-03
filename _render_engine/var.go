package RenderEngine

import (
	"bytes"
	"sync"
)

var (

	// Use a sync.Pool for bytes.Buffer to reduce allocations
	bufPool = sync.Pool{
		New: func() interface{} { return new(bytes.Buffer) },
	}
)
