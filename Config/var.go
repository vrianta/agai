package Config

// Config package provides configuration settings for the server

var (
	Http         = false
	Build        = false
	StaticFolder = []string{
		"Static",
	}
	CssFolder = []string{
		"Css",
	}
	JsFolders = []string{
		"Js",
	}
	ViewFolder = "Views"
)
