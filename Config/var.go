package Config

// Config package provides configuration settings for the server

var (
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
