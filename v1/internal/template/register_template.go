package template

import (
	"net/http"
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

				templateComponents[file_name] = createTemplateContext(folder_path, full_file_name, file_type)

				if file_name == "404" {
					http.HandleFunc("/404/", func(w http.ResponseWriter, r *http.Request) {
						t, _ := templateComponents[file_name]
						if !config.GetWebConfig().Build {
							// log.WriteLogf("Updating the Template")
							t.Update()
						}

						buf, _ := t.Execute("")
						w.Write(buf)

					})
				}
			} else { // register themes
				folder_name := object.Name()
				RegisterTheme(folder_name)
			}

		}
	}
}

func RegisterTheme(theme_folder string) {

	var found_404 bool = false // to see if the 404 page is present in the current theme or directory
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

				// templateRegistry[utils.JoinPath(theme_folder, file_name)] = createTemplateContext(file_full_path, full_file_name, file_type)
				templateComponents[theme_folder+"."+file_name] = createTemplateContext(file_full_path, full_file_name, file_type)

				if file_name == "404" {
					found_404 = true
					http.HandleFunc("/"+utils.JoinPath(theme_folder, file_name), func(w http.ResponseWriter, r *http.Request) {
						t, _ := templateComponents[file_name]

						if !config.GetWebConfig().Build {
							// log.WriteLogf("Updating the Template")
							t.Update()
						}

						buf, _ := t.Execute("")
						w.Write(buf)

					})
				}

			}
		}
	}

	if !found_404 {
		http.HandleFunc("/"+theme_folder+"/404/", func(w http.ResponseWriter, r *http.Request) {
			w.Write(_404__)
		})
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
