package Config

func New(
	Http bool,
	__Views_folder string,
	__Static_folders []string,
	__CSS_Folders []string,
	__JS_Folders []string,
) Class {

	config.Http = Http
	if __Static_folders != nil {
		config.Static_folders = __Static_folders
	}

	if __Views_folder != "" {
		config.Views_folder = __Views_folder
	}

	if __CSS_Folders != nil {
		config.CSS_Folders = __CSS_Folders
	}

	if __JS_Folders != nil {
		config.JS_Folders = __JS_Folders
	}

	return config
}

func Get() *Class {
	return &config
}
