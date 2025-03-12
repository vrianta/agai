package server

func newConfig(_config *Config) Config {

	if _config == nil {
		return config
	}

	config.Http = _config.Http
	if _config.Static_folders != nil {
		config.Static_folders = _config.Static_folders
	}

	if _config.Views_folder != "" {
		config.Views_folder = _config.Views_folder
	}

	return config
}
