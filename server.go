package Server

import (
	"net/http"

	"github.com/vrianta/Server/Config"
	"github.com/vrianta/Server/Log"
	"github.com/vrianta/Server/Router"
)

/*
 * server.New(hostname, port, routes, _config) -> function to create a instance of the server
 * @return -> it will return a pointer to the server with default
 * host -> is hostname of the server which host name you want to allow * is for everything and localhost to allow only local host connections
 * port -> the port number the server going to listen to
 * route ->  routes configaration which tells the
 * _config -> send the config of the server can be send nill if default is fine for you
 */
func New(host, port string, routes Routes) *_Struct {
	srvInstance = &_Struct{
		Host:   host,
		Port:   port,
		Router: Router.New(Router.Type(routes)),
	}
	return srvInstance
}

// Start runs the HTTP server
func (s *_Struct) Start() {

	Config.New() // Load the Configurations

	s.setup_static_folders()
	s.setup_css_folder()
	s.setup_js_folder()

	// setting up the Custom Routing Handler for the syste
	http.HandleFunc("/", s.Router.Handler)

	// Define the server configuration
	s.server = &http.Server{
		Addr: s.Host + ":" + s.Port, // Host and port
	}

	Log.WriteLogf("Server Starting at : %s:%s", s.Host, s.Port)

	s.server.ListenAndServe()
}

func (s *_Struct) setup_static_folders() {
	// Create a file server handler
	for static_folder := range Config.StaticFolders {
		fs := http.FileServer(http.Dir(Config.StaticFolders[static_folder]))
		Log.WriteLog("setting Up Static Folder : ", static_folder)
		http.Handle("/"+Config.StaticFolders[static_folder]+"/", http.StripPrefix("/"+Config.StaticFolders[static_folder]+"/", fs))
	}
}

// Generating Creating Routes for the Css Folders
func (s *_Struct) setup_css_folder() {
	for css_folder := range Config.CssFolders {
		http.HandleFunc("/"+Config.CssFolders[css_folder]+"/", s.Router.CSSHandlers)
		// s.Routes["/"+s.Config.CSS_Folders[css_folder]] = s.CSSHandlers
	}
}

func (s *_Struct) setup_js_folder() {
	for folder := range Config.JsFolders {
		http.HandleFunc("/"+Config.JsFolders[folder]+"/", s.Router.JsHandler)
		// s.Routes["/"+s.Config.CSS_Folders[folder]] = s.CSSHandlers
	}
}
