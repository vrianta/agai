package Config

// Config package provides configuration settings for the server

var (
	config = Class{
		Http: false,
		Static_folders: []string{
			"Static",
		},
		CSS_Folders: []string{
			"Css",
		},
		JS_Folders: []string{
			"Js",
		},
		Views_folder: "Views",
	}
)
