package template

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/vrianta/agai/v1/config"
	"github.com/vrianta/agai/v1/log"
	"github.com/vrianta/agai/v1/utils"
)

var view_folder = utils.JoinPath(".", config.GetViewFolder())

// each folder is a theme and it will be store as theme name in the template registry
func init() {

	if objects, err := os.ReadDir(view_folder); err != nil {
		log.Error("Failed to load read Directory %v", err)
		panic("Failed Loading View Folder")
	} else {
		for _, object := range objects {
			if !object.IsDir() {
				full_file_name := object.Name()                                    // name of file ex. hello.go
				file_type := strings.TrimPrefix(filepath.Ext(full_file_name), ".") // File extension/type
				file_name := full_file_name[:len(full_file_name)-len(file_type)-1] // Name without extension
				folder_path := utils.JoinPath(view_folder, full_file_name)
				//
				// c, err := registerTemplate(folder_path, folder_name)
				// if err != nil {
				// 	log.Error("Failed to Register Tempalte: %s - %s", folder_name, err)
				// 	panic("")
				// }
				templateRegistry[file_name] = createTemplateContext(folder_path, full_file_name, file_type)
			} else { // register themes
				folder_name := object.Name()
				RegisterTheme(folder_name)
				// c, err := RegisterTheme(folder_path, folder_name)
				// if err != nil {
				// 	log.Error("Failed to Register Tempalte: %s - %s", folder_name, err)
				// 	panic("")
				// }
				// templateRegistry[folder_name] = c.index
			}

		}
	}
}

func RegisterTheme(theme_folder string) {

	full_theme_path := utils.JoinPath(view_folder, theme_folder)

	if objects, err := os.ReadDir(full_theme_path); err != nil {
		log.Error("Failed to load read Directory %v", err)
		panic("Failed Loading View Folder")
	} else {
		for _, object := range objects {
			if object.IsDir() {
				RegisterTheme(utils.JoinPath(theme_folder, object.Name()))
			} else {
				full_file_name := object.Name()                                    // name of file ex. hello.go
				file_type := strings.TrimPrefix(filepath.Ext(full_file_name), ".") // File extension/type
				file_name := full_file_name[:len(full_file_name)-len(file_type)-1] // Name without extension

				file_full_path := utils.JoinPath(full_theme_path, full_file_name)

				templateRegistry[utils.JoinPath(theme_folder, file_name)] = createTemplateContext(file_full_path, full_file_name, file_type)
				templateComponents[theme_folder+"."+file_name] = createTemplateContext(file_full_path, full_file_name, file_type)
			}
		}
	}
}

func createTemplateContext(view_path, full_file_name, file_type string) *Context {

	// Register the template using the custom Template package
	if _template, err := create(view_path, full_file_name, file_type, config.GetBuild()); err != nil {
		log.Error("Failed to Create the template: %s Error: %v", view_path, err)
		panic("")
	} else {
		return _template
	}

}
