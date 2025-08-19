package server

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/vrianta/agai/v1/config"
	"github.com/vrianta/agai/v1/database"
	Session "github.com/vrianta/agai/v1/internal/session"
	Log "github.com/vrianta/agai/v1/log"
	"github.com/vrianta/agai/v1/model"
	Router "github.com/vrianta/agai/v1/router"
)

/*
 * server.New(hostname, port, routes, _config) -> function to create a instance of the server
 * @return -> it will return a pointer to the server with default
 * host -> is hostname of the server which host name you want to allow * is for everything and localhost to allow only local host connections
 * port -> the port number the server going to listen to
 * route ->  routes configaration which tells the
 * _config -> send the config of the server can be send nill if default is fine for you
 */

// waitForPort waits until the port is available or timeout occurs
func waitForPort(addr string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		ln, err := net.Listen("tcp", addr)
		if err == nil {
			ln.Close()
			return nil // port is free
		}
		Log.WriteLogf("[WARN]: Port %s is busy, waiting...\n", addr)
		time.Sleep(500 * time.Millisecond)
	}
	return fmt.Errorf("port %s did not become free in time", addr)
}

// Start runs the HTTP server
func Start() *instance {

	s := &instance{}

	database.Init()
	s.setup_static_folders()
	s.setup_css_folder()
	s.setup_js_folder()
	s.setup_views() // Register all the views with the RenderEngine

	model.Init() // intialsing model with creating tables and updating them

	// Initialize Models Handler

	// Starting Session Handler to Manage Session Expiry
	go Session.StartSessionHandler()
	go Session.StartLRUHandler() // Start the LRU handler for session management

	// setting up the Custom Routing Handler for the syste
	http.HandleFunc("/", routerHandler)

	// Define the server configuration
	addr := config.GetWebConfig().Host + ":" + config.GetWebConfig().Port
	s.server = &http.Server{
		Addr: addr,
	}

	// Wait for port to be free before starting
	if err := waitForPort(addr, 20*time.Second); err != nil {
		panic("[Server] " + err.Error())
	}

	if config.StartServer {
		Log.WriteLogf("[Server] Started at : http://localhost:%s\n", config.GetWebConfig().Port)
		fmt.Print("---------------------------------------------------------\n\n")

		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic("[Server] Failed to start: " + err.Error())
		}
	}

	return s
}

// if the user want to have a custom Routing Handler the they can use this function to register it
func (s *instance) RegisterCustomRoutingHandler(_func func(w http.ResponseWriter, r *http.Request)) {
	routerHandler = _func
}

func (s *instance) setup_static_folders() {
	// Create a file server handler
	for _, folder := range config.GetWebConfig().StaticFolders {
		fs := http.FileServer(http.Dir(folder))
		http.Handle("/"+folder+"/", http.StripPrefix("/"+folder+"/", fs))
	}
}

// Generating Creating Routes for the Css Folders
func (s *instance) setup_css_folder() {
	for _, folder := range config.GetWebConfig().CssFolders {
		http.HandleFunc("/"+folder+"/", Router.StaticFileHandler("text/css; charset=utf-8"))
	}
}

func (s *instance) setup_js_folder() {
	for _, folder := range config.GetWebConfig().JsFolders {
		http.HandleFunc("/"+folder+"/", Router.StaticFileHandler("application/javascript; charset=utf-8"))
	}
}

// function to go through all the routes and register their Views and create templates
func (s *instance) setup_views() {

	routes := Router.GetRoutes()
	fmt.Print("---------------------------------------------------------\n")
	fmt.Print("[Views Setup] Registering templates for controllers:\n")
	fmt.Print("---------------------------------------------------------\n")
	for route, controller := range *routes {
		if controller.View != "" {
			if err := controller.RegisterTemplate(); err != nil {
				fmt.Printf("[Error]   Template: %-20s | %v\n", controller.View, err)
				os.Exit(1)
			} else {
				fmt.Printf("[Success] Template: %s and Route: %-20s | Registered\n", controller.View, route)
			}
		}
	}
	fmt.Print("---------------------------------------------------------\n")
	fmt.Print("[Views Setup] All views registered successfully.\n")
	fmt.Print("---------------------------------------------------------\n\n")
}
