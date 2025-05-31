package Server

import (
	"net/http"

	"github.com/vrianta/Server/Router"
)

type (
	Routes Router.Type
	// Server represents the HTTP server with session management
	_Struct struct {
		Router *Router.Struct
		server *http.Server

		state bool // hold 0 or 1 to ensure if the server is runnning or not
	}
)
