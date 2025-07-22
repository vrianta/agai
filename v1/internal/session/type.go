package session

import (
	"time"
)

type (
	SessionHeap []*Instance
	SessionVars map[string]any
	PostParams  map[string]string
	GetParams   map[string]string

	Instance struct {
		ID string

		POST  PostParams
		GET   GetParams
		Store SessionVars

		LoggedIn bool

		Expiry time.Time // Expiry time for the session
		// lastUsed atomic.Int64
	}

	lruUpdate struct {
		SessionID string
		Op        string // "move" or "InsertRow"
	}
)
