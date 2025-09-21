package server

import "github.com/vrianta/agai/v1/internal/requestHandler"

// Global instance of the server
var (
	routerHandler = requestHandler.Handler
)
