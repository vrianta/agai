package server

// Global instance of the server
var (
	srvInstance *server

	config = Config{
		Http:          false,
		Static_folder: "Static",
		Views_folder:  "Views",
	}
)
