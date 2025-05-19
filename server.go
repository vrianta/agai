package server

import (
	"net/http"
)

/*
 * server.New(hostname, port, routes, _config) -> function to create a instance of the server
 * @return -> it will return a pointer to the server with default
 * host -> is hostname of the server which host name you want to allow * is for everything and localhost to allow only local host connections
 * port -> the port number the server going to listen to
 * route ->  routes configaration which tells the
 * _config -> send the config of the server can be send nill if default is fine for you
 */
func New(host, port string, routes Routes, _config *Config) *server {
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

	s.setup_static_folders()
	s.setup_css_folder()
	s.setup_js_folder()

	// setting up the Custom Routing Handler for the syste
	http.HandleFunc("/", s.routingHandler)

	// Define the server configuration
	s.server = &http.Server{
		Addr: s.Host + ":" + s.Port, // Host and port
	}

	WriteLogf("Server Starting at : %s:%s", s.Host, s.Port)

	s.server.ListenAndServe()
	// s.state = true
	// s.ServeConsole()

	// for s.state {
	// }

}

func (s *server) setup_static_folders() {
	// Create a file server handler
	for static_folder := range s.Config.Static_folders {
		fs := http.FileServer(http.Dir(s.Config.Static_folders[static_folder]))
		WriteLog("setting Up Static Folder : ", static_folder)
		http.Handle("/"+s.Config.Static_folders[static_folder]+"/", http.StripPrefix("/"+s.Config.Static_folders[static_folder]+"/", fs))
	}
}

// Generating Creating Routes for the Css Folders
func (s *server) setup_css_folder() {
	for css_folder := range s.Config.CSS_Folders {
		http.HandleFunc("/"+s.Config.CSS_Folders[css_folder]+"/", s.CSSHandlers)
		// s.Routes["/"+s.Config.CSS_Folders[css_folder]] = s.CSSHandlers
	}
}

func (s *server) setup_js_folder() {
	for folder := range s.Config.JS_Folders {
		http.HandleFunc("/"+s.Config.JS_Folders[folder]+"/", s.JsHandler)
		// s.Routes["/"+s.Config.CSS_Folders[folder]] = s.CSSHandlers
	}
}

// RemoveSession removes a session from the session manager
func RemoveSession(sessionID string) {
	defer delete(srvInstance.Sessions, sessionID)
}
