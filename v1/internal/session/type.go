package session

import (
	"time"
)

type (
	// Collection represents a collection of sessions managed as a heap
	collection []*Instance

	// SessionData stores arbitrary key-value data for a session
	SessionData map[string]any

	// HTTPPostParameters represents POST request parameters
	HTTPPostParameters map[string]string

	// HTTPGetParameters represents GET request parameters
	HTTPGetParameters map[string]string

	Instance struct {
		ID string

		PostParameters HTTPPostParameters
		GetParameters  HTTPGetParameters
		Data           SessionData

		IsAuthenticated bool

		ExpirationTime time.Time
		// lastUsed atomic.Int64
	}

	// LRUCacheOperation represents an operation to update the LRU cache
	LRUCacheOperation struct {
		ID            string
		OperationType string // "move_to_front" or "insert" or "remove"
	}
)
