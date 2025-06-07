package Session

import (
	"net/http"
	"time"

	"github.com/vrianta/Server/RenderEngine"
)

type (
	SessionVars map[string]any
	PostParams  map[string]string
	GetParams   map[string]string

	Struct struct {
		ID string
		W  http.ResponseWriter
		R  *http.Request

		POST  PostParams
		GET   GetParams
		Store SessionVars

		RenderEngine RenderEngine.Struct

		Expiry   time.Time // Expiry time for the session
		LastUsed time.Time // Last access time for LRU
	}

	lruUpdate struct {
		SessionID string
		Op        string // "move" or "insert"
	}
)
