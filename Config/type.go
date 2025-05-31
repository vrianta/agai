package Config

// Config package provides configuration settings for the server

// Class represents the configuration for the server
// It includes settings for HTTP/HTTPS, static files, CSS, JS, and views folders.
type (

	// Flaggs for Server Config where it will care of the config of the server
	// http ->  is to tell server if it need to load https or http server for example http enabled mean it will load http server else by default it will be https
	// By Default the Static Files will be in /Static and can be accessed in html by Static/files_path
	Class struct {
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
)
