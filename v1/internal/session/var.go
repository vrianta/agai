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

// count of total sessions are there in the system right now
var createdSessionCount int

// type SessionStoreType string

const (
	session_store_type_memory   string = "memory"
	session_store_type_disk     string = "disk"
	session_store_type_storage  string = "storage"
	session_store_type_db       string = "db"
	session_store_type_database string = "database"
	// Add more types as needed, for example:
	// Disk
	// Redis
)
