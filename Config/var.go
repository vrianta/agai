package Config

// Config package provides configuration settings for the server

var (
	webConfigFile = "config.web.json"
	dBConfigFile  = "config.database.json"
	webConfig     = webConfigStruct{
		Port: "1080", // Default port for the server
		Host: "",     // Default host for the server

		// Flaggs for Server Config where it will care of the config of the server
		Https: false,
		Build: false,
		StaticFolders: []string{
			"Static",
		},
		CssFolders: []string{
			"Css",
		},
		JsFolders: []string{
			"Js",
		},
		ViewFolder:       "Views",
		MaxSessionCount:  50000,    // Default value.
		SessionStoreType: "memory", // Default session store type
	}

	databaseConfig = databaseConfigStruct{
		Host:     "",
		Port:     "3306", // Default MySQL port
		User:     "root",
		Password: "",
		Database: "mydatabase",
		Protocol: "tcp", // Default protocol for MySQL
	}
)
