package server

import (
	"net/http"
)

type (
	// Server represents the HTTP server with session management
	instance struct {
		server *http.Server

		state bool // hold 0 or 1 to ensure if the server is runnning or not
	}
)
