package session

import (
	"container/list"
	"sync"
)

var (
	all               = make(map[string]*Struct)       // Session map
	lruList           = list.New()                     // Doubly-linked list for LRU
	lruMap            = make(map[string]*list.Element) // Map session ID to list element
	mutex             = sync.RWMutex{}
	sessionWakeupChan = make(chan struct{}, 1)
	lruUpdateChan     = make(chan lruUpdate, 1000) // Buffered channel for LRU ops

	updateMutex = sync.Mutex{}
	cleanMutex  = sync.Mutex{}
)
