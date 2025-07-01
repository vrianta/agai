package Server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/vrianta/Server/Config"
	"github.com/vrianta/Server/DatabaseHandler"
	"github.com/vrianta/Server/Log"
	"github.com/vrianta/Server/Router"
	"github.com/vrianta/Server/Session"
	Models "github.com/vrianta/Server/modelsHandler"
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
	srvInstance = &_Struct{}
	return srvInstance
}

// Start runs the HTTP server
func (s *_Struct) Start() *_Struct {

	s.setup_static_folders()
	s.setup_css_folder()
	s.setup_js_folder()
	s.setup_views()      // Register all the views with the RenderEngine
	s.initialiseModels() // intialsing models with creating tables and updating them

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

	Log.WriteLogf("[Server] Started at : http://localhost:%s\n", Config.GetWebConfig().Port)

	if err := s.server.ListenAndServe(); err != nil {
		panic("[Server] Failed to start: " + err.Error())
	}

	return s
}

// Initialise Database Handler
func (s *_Struct) InitDatabase() *_Struct {

	if err := DatabaseHandler.Init(); err != nil {
		panic("Database Initialisation failed: " + err.Error())
	} else {
		Log.WriteLog("[Database] Initialised Successfully\n")
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
	fmt.Print("---------------------------------------------------------\n")
	fmt.Print("[Views Setup] Registering templates for controllers:\n")
	fmt.Print("---------------------------------------------------------\n")
	for _, controller := range *routes {
		if controller.View != "" {
			if err := controller.RegisterTemplate(); err != nil {
				fmt.Printf("[Error]   Template: %-20s | %v\n", controller.View, err)
				os.Exit(1)
			} else {
				fmt.Printf("[Success] Template: %-20s | Registered\n", controller.View)
			}
		}
	}
	fmt.Print("---------------------------------------------------------\n")
	fmt.Print("[Views Setup] All views registered successfully.\n")
	fmt.Print("---------------------------------------------------------\n\n")
}

func (s *_Struct) initialiseModels() {
	if Config.GetBuild() {
		fmt.Print("[Models] Build mode enabled, skipping model initialization.\n")
		return
	}
	fmt.Print("---------------------------------------------------------\n")
	fmt.Print("[Models] Initializing models and syncing database tables:\n")
	fmt.Print("---------------------------------------------------------\n")
	for _, model := range Models.ModelsRegistry {
		if DatabaseHandler.Initialized {
			fmt.Printf("[Model]   Table: %-20s | Syncing...\n", model.TableName)
			model.GetTableScema()
			model.CreateTableIfNotExists() // creating table if not existed
		} else {
			fmt.Printf("[Warning] Database not initialized, skipping table creation for model: %-20s\n", model.TableName)
		}
	}
	fmt.Print("---------------------------------------------------------\n")
	fmt.Print("[Models] Model initialization complete.\n")
	fmt.Print("---------------------------------------------------------\n\n")

}
