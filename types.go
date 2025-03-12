package server

import (
	"html/template"
	"net/http"
	"time"
)

type (
	ResponseCode int // server response codes
	Uri          string

	// RoutesMap is a type alias for mapping routes to handlers
	RoutesMap   map[string]func(*Session)
	SessionVars map[string]interface{}
	PostParams  map[string]string
	GetParams   map[string]string

	// Server represents the HTTP server with session management
	server struct {
		Host     string
		Port     string
		Routes   RoutesMap
		Config   Config
		Sessions map[string]Session

		server *http.Server

		state bool // hold 0 or 1 to ensure if the server is runnning or not
	}

	templates struct {
		Uri          string            // path of the template file
		LastModified time.Time         // date when the file last modified
		Data         template.Template // template data of the file before modified
	}

	RenderEngine struct {
		view []byte
		W    http.ResponseWriter
	}

	RenderData map[string]interface{}

	// Flaggs for Server Config where it will care of the config of the server
	// http ->  is to tell server if it need to load https or http server for example http enabled mean it will load http server else by default it will be https
	// By Default the Static Files will be in /Static and can be accessed in html by Static/files_path
	Config struct {
		Http           bool
		Static_folders []string // static folders list which comes with a list of scrigs for all the static folder which needs to be in file server
		// list of folder where the CSS files will be keep that in mind that Static folder also can have
		// the Css files but image sure keep Css file only in CSS folders is good idea because it make this system faster no randome file checks needed
		// Remember this CSS file will be auto loaded and will be updated once you update them in the local better for development and build system
		// when we will introduce build systems
		CSS_Folders  []string
		JS_Folders   []string
		Views_folder string
	}

	Session struct {
		ID string
		w  http.ResponseWriter
		r  *http.Request

		POST  PostParams
		GET   GetParams
		Store SessionVars

		RenderEngine RenderEngine
	}
)
