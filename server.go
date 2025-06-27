package Server

import (
	"net/http"
	"os"

	"github.com/vrianta/Server/Config"
	"github.com/vrianta/Server/DatabaseHandler"
	"github.com/vrianta/Server/Log"
	"github.com/vrianta/Server/Router"
	"github.com/vrianta/Server/Session"

	"database/sql"
)

/*
 * server.New(hostname, port, routes, _config) -> function to create a instance of the server
 * @return -> it will return a pointer to the server with default
 * host -> is hostname of the server which host name you want to allow * is for everything and localhost to allow only local host connections
 * port -> the port number the server going to listen to
 * route ->  routes configaration which tells the
 * _config -> send the config of the server can be send nill if default is fine for you
 */
func New() *_Struct {
	Config.Init() // Load the Configurations

	srvInstance = &_Struct{}
	return srvInstance
}

// Start runs the HTTP server
func (s *_Struct) Start() *_Struct {

	s.setup_static_folders()
	s.setup_css_folder()
	s.setup_js_folder()
	s.setup_views() // Register all the views with the RenderEngine

	// Initialize Models Handler

	// Starting Session Handler to Manage Session Expiry
	go Session.StartSessionHandler()
	go Session.StartLRUHandler() // Start the LRU handler for session management

	// setting up the Custom Routing Handler for the syste
	http.HandleFunc("/", Router.Handler)

	// Define the server configuration
	s.server = &http.Server{
		Addr: Config.GetWebConfig().Host + ":" + Config.GetWebConfig().Port, // Host and port
	}

	Log.WriteLogf("Server Starting at : http://localhost:%s\n", Config.GetWebConfig().Port)

	if err := s.server.ListenAndServe(); err != nil {
		panic("Server failed to start: " + err.Error())
	}

	return s
}

// Initialise Database Handler
func (s *_Struct) RegisterDatabase(sql *sql.DB) *_Struct {

	if err := DatabaseHandler.Init(sql); err != nil {
		panic("Database Initialisation failed: " + err.Error())
	} else {
		Log.WriteLog("Database Initialised Successfully")
	}

	return s
}

func (s *_Struct) setup_static_folders() {
	// Create a file server handler
	for _, folder := range Config.GetWebConfig().StaticFolders {
		fs := http.FileServer(http.Dir(folder))
		http.Handle("/"+folder+"/", http.StripPrefix("/"+folder+"/", fs))
	}
}

// Generating Creating Routes for the Css Folders
func (s *_Struct) setup_css_folder() {
	for _, folder := range Config.GetWebConfig().CssFolders {
		http.HandleFunc("/"+folder+"/", Router.StaticFileHandler("text/css; charset=utf-8"))
	}
}

func (s *_Struct) setup_js_folder() {
	for _, folder := range Config.GetWebConfig().JsFolders {
		http.HandleFunc("/"+folder+"/", Router.StaticFileHandler("application/javascript; charset=utf-8"))
	}
}

// function to go through all the routes and register their Views and create templates
func (s *_Struct) setup_views() {

	routes := Router.GetRoutes()
	for _, controller := range *routes {
		if controller.View != "" {
			if err := controller.RegisterTemplate(); err != nil {
				Log.WriteLogf("Error registering template %s: %v", controller.View, err)
				os.Exit(1)
			} else {
				Log.WriteLogf("Template registered: %s\n", controller.View)
			} // Register the template with the RenderEngine
		}
	}
	Log.WriteLog("Views setup completed")
}
