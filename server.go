package server

import (
	"net/http"
)

// RoutesMap is a type alias for mapping routes to handlers
type RoutesMap map[string]func(*Session)

// Server represents the HTTP server with session management
type server struct {
	Host     string
	Port     string
	Routes   RoutesMap
	Config   Config
	Sessions map[string]Session

	server *http.Server

	state bool // hold 0 or 1 to ensure if the server is runnning or not
}

// Global instance of the server
var (
	srvInstance *server
)

/*
 * server.New(hostname, port, routes, _config) -> function to create a instance of the server
 * @return -> it will return a pointer to the server with default
 * host -> is hostname of the server which host name you want to allow * is for everything and localhost to allow only local host connections
 * port -> the port number the server going to listen to
 * route ->  routes configaration which tells the
 * _config -> send the config of the server can be send nill if default is fine for you
 */
func New(host, port string, routes RoutesMap, _config *Config) *server {
	srvInstance = &server{
		Host:     host,
		Port:     port,
		Routes:   routes,
		Config:   newConfig(_config),
		Sessions: make(map[string]Session),
	}
	return srvInstance
}

// Start runs the HTTP server
func (s *server) Start() {

	// Create a file server handler
	fs := http.FileServer(http.Dir(s.Config.Static_folder))

	// Serve static files from the /static/ URL path
	http.Handle("/Static/", http.StripPrefix("/Static/", fs))

	http.HandleFunc("/", s.routingHandler)

	// Define the server configuration
	s.server = &http.Server{
		Addr: s.Host + ":" + s.Port, // Host and port
	}

	WriteLogf("Server Starting at : %s:%s", s.Host, s.Port)

	go s.server.ListenAndServe()
	s.state = true
	s.ServeConsole()

	// s.server.

}

// RemoveSession removes a session from the session manager
func RemoveSession(sessionID string) {
	delete(srvInstance.Sessions, sessionID)
}
