package Config

// Config package provides configuration settings for the server

var (
	Port string = "1080" // Default port for the server
	Host string = ""     // Default host for the server

	// Flaggs for Server Config where it will care of the config of the server
	Http          = false
	Build         = false
	StaticFolders = []string{
		"Static",
	}
	CssFolders = []string{
		"Css",
	}
	JsFolders = []string{
		"Js",
	}
	ViewFolder = "Views"
)
