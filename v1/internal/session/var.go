package session

import (
	"container/list"
	"sync"
)

var (
	// Core session storage
	instances = make(map[string]*Instance) // Maps session ID to session instance

	lruOrderList  = list.New()                     // Doubly-linked list maintaining LRU order
	lruElementMap = make(map[string]*list.Element) // Maps session ID to its position in LRU list

	// Channel-based communication
	cleanupTriggerChan = make(chan struct{}, 1) // Triggers cleanup goroutine

	lruOperationChan = make(chan LRUCacheOperation, 1000) // Buffered channel for LRU operations

	// Synchronization primitives
	sessionUpdateMutex = sync.Mutex{} // Protects session creation/updates
	sessionCleanMutex  = sync.Mutex{} // Protects cleanup operations

	sessionStoreMutex = sync.RWMutex{} // Protects session storage map
	lruCacheMutex     = sync.RWMutex{} // Protects LRU cache operations
)

// sessionHeap
var (
	sessionHeap     collection
	heapAccessMutex sync.Mutex
)
