package server

import (
	"fmt"
	"net/http"
)

// RoutesMap is a type alias for mapping routes to handlers
type RoutesMap map[string]func(*Session)

// Server represents the HTTP server with session management
type Server struct {
	Host     string
	Port     string
	Routes   RoutesMap
	Sessions map[string]Session
}

// Global instance of the server
var srvInstance *Server

// NewServer creates a new instance of the Server
func New(host, port string, routes RoutesMap) *Server {
	srvInstance = &Server{
		Host:     host,
		Port:     port,
		Routes:   routes,
		Sessions: make(map[string]Session),
	}
	return srvInstance
}

// Start runs the HTTP server
func (s *Server) Start() error {
	http.HandleFunc("/", s.routingHandler)

	// Define the server configuration
	server := &http.Server{
		Addr: s.Host + ":" + s.Port, // Host and port
	}

	WriteConsole("Server Starting at " + s.Host + ":" + s.Port)

	// Start the server
	if err := server.ListenAndServe(); err != nil {
		fmt.Println("Error starting server:", err)
		return err
	}

	return nil
}

// RemoveSession removes a session from the session manager
func RemoveSession(sessionID string) {
	delete(srvInstance.Sessions, sessionID)
}
