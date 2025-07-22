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

		// PostParameters HTTPPostParameters
		// GetParameters  HTTPGetParameters
		Data SessionData

		IsAuthenticated bool

		ExpirationTime time.Time
		// lastUsed atomic.Int64
		
		heapIndex int // Index in the session heap for efficient updates
	}

	// LRUCacheOperation represents an operation to update the LRU cache
	LRUCacheOperation struct {
		ID            string
		OperationType string // "move_to_front" or "insert" or "remove"
	}
)

// Current: Push to heap on every store
// heap.Push(&sessionHeap, session)

// Improvement: Only update heap when expiration changes
// if session.ExpirationTime != oldExpiration {
//     heap.Fix(&sessionHeap, session.heapIndex)
// }
