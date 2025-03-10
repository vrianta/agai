package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// RoutesMap is a type alias for mapping routes to handlers
type RoutesMap map[string]func(*Session)

// Server represents the HTTP server with session management
type Server struct {
	Host     string
	Port     string
	Routes   RoutesMap
	Sessions map[string]Session

	server *http.Server
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
func (s *Server) Start() {
	http.HandleFunc("/", s.routingHandler)

	// Define the server configuration
	s.server = &http.Server{
		Addr: s.Host + ":" + s.Port, // Host and port
	}

	WriteLogf("Server Starting at " + s.Host + ":" + s.Port)

	go s.server.ListenAndServe()
	s.ServeConsole()

	// s.server.

}

func (s *Server) ServeConsole() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt) // Listen for Ctrl+C

	go func() {
		<-quit // Wait for interrupt signal
		fmt.Println("\nShutting down server...")
		s.stopServer()
		os.Exit(0) // Exit program gracefully
	}()

	for {
		var input string
		fmt.Print(": ")
		fmt.Scanln(&input)

		switch input {
		case "stop":
			s.stopServer()
		case "start":
			s.startServer()
		case "exit":
			fmt.Println("Exiting...")
			os.Exit(0)
		}
	}
}

func (s *Server) startServer() {
	// Define the server configuration
	s.server = &http.Server{
		Addr: s.Host + ":" + s.Port, // Host and port
	}

	WriteLogf("Server Starting at " + s.Host + ":" + s.Port)

	go s.server.ListenAndServe()
}

func (s *Server) stopServer() {
	// Create a timeout context (5 seconds)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown server gracefully
	if err := s.server.Shutdown(ctx); err != nil {
		fmt.Println("Shutdown Error:", err)
	} else {
		fmt.Println("Server shutdown successfully")
	}
}

// RemoveSession removes a session from the session manager
func RemoveSession(sessionID string) {
	delete(srvInstance.Sessions, sessionID)
}
