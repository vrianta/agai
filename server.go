package Server

import (
	"net/http"
	"os"

	"github.com/vrianta/Server/Config"
	"github.com/vrianta/Server/Log"
	"github.com/vrianta/Server/RenderEngine"
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
func New(routes Routes) *_Struct {
	Config.Init() // Load the Configurations

	srvInstance = &_Struct{
		Router: Router.New(Router.Type(routes)),
	}
	return srvInstance
}

// Start runs the HTTP server
func (s *_Struct) Start() {

	s.setup_static_folders()
	s.setup_css_folder()
	s.setup_js_folder()
	s.setup_views() // Register all the views with the RenderEngine

	// setting up the Custom Routing Handler for the syste
	http.HandleFunc("/", s.Router.Handler)

	// Define the server configuration
	s.server = &http.Server{
		Addr: Config.Host + ":" + Config.Port, // Host and port
	}

	Log.WriteLogf("Server Starting at : http://%s:%s", Config.Host, Config.Port)

	s.server.ListenAndServe()
}

func (s *_Struct) setup_static_folders() {
	// Create a file server handler
	for _, folder := range Config.StaticFolders {
		fs := http.FileServer(http.Dir(folder))
		http.Handle("/"+folder+"/", http.StripPrefix("/"+folder+"/", fs))
	}
}

// Generating Creating Routes for the Css Folders
func (s *_Struct) setup_css_folder() {
	for _, folder := range Config.CssFolders {
		http.HandleFunc("/"+folder+"/", s.Router.StaticFileHandler("text/css; charset=utf-8"))
	}
}

func (s *_Struct) setup_js_folder() {
	for _, folder := range Config.JsFolders {
		http.HandleFunc("/"+folder+"/", s.Router.StaticFileHandler("application/javascript; charset=utf-8"))
	}
}

// function to go through all the routes and register their Views and create templates
func (s *_Struct) setup_views() {
	routes := s.Router.Get()
	for _, route := range *routes {
		if route.View != "" {
			if err := RenderEngine.RegisterTemplate(route.View); err != nil {
				Log.WriteLogf("Error registering template %s: %v", route.View, err)
				os.Exit(1)
			} else {
				Log.WriteLogf("Template registered: %s", route.View)
			} // Register the template with the RenderEngine
		}
	}
	Log.WriteLog("Views setup completed")
}
