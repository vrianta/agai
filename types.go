package Server

import (
	"net/http"

	Config "github.com/vrianta/Server/Config"
	"github.com/vrianta/Server/Router"
)

type (
	Routes Router.Type
	// Server represents the HTTP server with session management
	_Struct struct {
		Host   string
		Port   string
		Router *Router.Struct
		Config Config.Class

		server *http.Server

		state bool // hold 0 or 1 to ensure if the server is runnning or not
	}
)
