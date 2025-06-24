package Session

import (
	"time"
)

type (
	SessionVars map[string]any
	PostParams  map[string]string
	GetParams   map[string]string

	Struct struct {
		ID string

		POST  PostParams
		GET   GetParams
		Store SessionVars

		isLoggedIn bool

		Expiry   time.Time // Expiry time for the session
		LastUsed time.Time // Last access time for LRU
	}

	lruUpdate struct {
		SessionID string
		Op        string // "move" or "insert"
	}
)
