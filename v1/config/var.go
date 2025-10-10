package config

// Config package provides configuration settings for the server

var (
	webConfigFile = "appsettings.json"
	dBConfigFile  = "config.database.json"

	webConfig = webConfigStruct{
		Port: "1080", // Default port for the server
		Host: "",     // Default host for the server

		// Flags for Server Config
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
		ViewFolder:       "views",
		MaxSessionCount:  50000,     // Default value
		SessionStoreType: "storage", // Default session store type
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
