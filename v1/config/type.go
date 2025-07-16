package config

// Config package provides configuration settings for the server

// Class represents the configuration for the server
// It includes settings for HTTP/HTTPS, static files, CSS, JS, and views folders.
type (

	// Flaggs for Server Config where it will care of the config of the server
	// http ->  is to tell server if it need to load https or http server for example http enabled mean it will load http server else by default it will be https
	// By Default the Static Files will be in /Static and can be accessed in html by Static/files_path
	webConfigStruct struct {
		Port          string   `json:"Port"`
		Host          string   `json:"Host"`
		Https         bool     `json:"Https"`
		Build         bool     `json:"Build"`
		StaticFolders []string `json:"StaticFolders"`
		// list of folder where the CSS files will be keep that in mind that Static folder also can have
		// the Css files but image sure keep Css file only in CSS folders is good idea because it make this system faster no randome file checks needed
		// Remember this CSS file will be auto loaded and will be updated once you update them in the local better for development and build system
		// when we will introduce build systems
		CssFolders       []string `json:"CssFolders"`       // css folders list which comes with a list of scrigs for all the css folder which needs to be in file server
		JsFolders        []string `json:"JsFolders"`        // js folders list which comes with a list of scrigs for all the js folder which needs to be in file server
		ViewFolder       string   `json:"ViewFolder"`       // view folder is the folder where all the views will be kept and it will be used to render the views
		MaxSessionCount  int      `json:"MaxSessionCount"`  // <-- Add this line
		SessionStoreType string   `json:"SessionStoreType"` // Type of session store, e.g., "memory", "redis", "database".
	}

	databaseConfigStruct struct {
		Host     string `json:"Host"`
		Port     string `json:"Port"`
		User     string `json:"User"`
		Password string `json:"Password"`
		Database string `json:"Database"`
		Protocol string `json:"Protocol"` // e.g., "tcp", "unix", etc.
		Driver   string `json:"Driver"`   // e.g., "mysql", "postgres", etc.
		SSLMode  string `json:"SSLMode"`  // e.g., "disable", "require", etc.
	}

	smtpConfigStruct struct {
		Host     string `json:"Host"`     // SMTP server host (e.g., smtp.gmail.com)
		Port     int    `json:"Port"`     // SMTP server port (usually 587 or 465)
		Username string `json:"Username"` // SMTP username (email address)
		Password string `json:"Password"` // SMTP password or app password
		UseTLS   bool   `json:"UseTLS"`   // Whether to use TLS encryption (default: true)
	}
)
