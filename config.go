package server

// Flaggs for Server Config where it will care of the config of the server
// http ->  is to tell server if it need to load https or http server for example http enabled mean it will load http server else by default it will be https
// By Default the Static Files will be in /Static and can be accessed in html by Static/files_path
type Config struct {
	Http          bool
	Static_folder string
	Views_folder  string
}

var config = Config{
	Http:          false,
	Static_folder: "Static",
	Views_folder:  "Views",
}

func newConfig(_config *Config) Config {

	if _config == nil {
		return config
	}

	config.Http = _config.Http
	if _config.Static_folder != "" {
		config.Static_folder = _config.Static_folder
	}

	if _config.Views_folder != "" {
		config.Views_folder = _config.Views_folder
	}

	return config
}
