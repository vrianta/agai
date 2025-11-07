package template

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/vrianta/agai/v1/config"
	"github.com/vrianta/agai/v1/log"
	"github.com/vrianta/agai/v1/utils"
)

var view_folder = utils.JoinPath(".", config.GetViewFolder())

// func init() {

// 	if folders, err := os.ReadDir(view_folder); err != nil {
// 		log.Error("Failed to load read Directory %v", err)
// 		panic("Failed Loading View Folder")
// 	} else {
// 		for _, folder := range folders {
// 			folder_name := folder.Name()

// 			if !folder.IsDir() {
// 				log.Warn("Files are not Reccomented in the View Directory - %s is found in %s", folder_name, view_folder)
// 				continue
// 			}

// 			folder_path := utils.JoinPath(view_folder, folder_name)
// 			c, err := registerTemplate(folder_path, folder_name)
// 			if err != nil {
// 				log.Error("Failed to Register Tempalte: %s - %s", folder_name, err)
// 				panic("")
// 			}
// 			templateRegistry[folder_name] = c
// 		}
// 	}

// }

/*
RegisterTemplate scans the controller's view directory and registers templates for each HTTP method.
It expects files named default.html/php/gohtml, get.html/php, post.html/php, etc.
Panics if no default view is found.
@param - view_path is path of the view folder where the templates are present
@return - Contexts of the view
Returns:
- error: if reading the directory or registering a template fails.
*/
func registerTemplate(view_path string, view_name string) (*Contexts, error) {
	// view name is the folder name where the view files are present
	// view_path is the exact full path of the view

	// fmt.Printf("Registering templates for controller: %T, view path: %s\n", c, view_path)
	files, err := os.ReadDir(view_path)
	if err != nil {
		err := fmt.Errorf("error reading directory: %s", err.Error())
		panic(err)
	}

	c := Contexts{}
	var gotDefaultView = false // Track if a default view is found
	for _, entry := range files {
		if !entry.IsDir() {
			full_file_name := entry.Name()                                        // name of file ex. hello.go
			var file_type = strings.TrimPrefix(filepath.Ext(full_file_name), ".") // File extension/type
			file_name := full_file_name[:len(full_file_name)-len(file_type)-1]    // Name without extension

			// Register the template using the custom Template package
			if _template, err := create(view_path, full_file_name, file_type); err != nil {
				println("testing")
				return nil, err
			} else {
				// fmt.Printf("  Found template: %s (type: %s) for controller: %T and file_name:%s Path:%s\n", full_file_name, file_type, c, file_name, view_path)
				switch file_name {
				case "default", "index":
					if c.index != nil {
						log.Error("You have two Defailt templates for %s, which is not allowed", view_path)
						panic("")
					}
					c.index = _template
					gotDefaultView = true
				case "get":
					if c.get != nil {
						log.Error("You have two Get templates for %s, which is not allowed", view_path)
						panic("")
					}
					c.get = _template
				case "post":
					if c.post != nil {
						log.Error("You have two Post templates for %s, which is not allowed", view_path)
						panic("")
					}
					c.post = _template
				case "delete":
					if c.delete != nil {
						log.Error("You have two Delete templates for %s, which is not allowed", view_path)
						panic("")
					}
					c.delete = _template
				case "patch":
					if c.patch != nil {
						log.Error("You have two Patch templates for %s, which is not allowed", view_path)
						panic("")
					}
					c.patch = _template
				case "put":
					if c.put != nil {
						log.Error("You have two PUT templates for %s, which is not allowed", view_path)
						panic("")
					}
					c.put = _template
				case "head":
					if c.head != nil {
						log.Error("You have two Head templates for %s, which is not allowed", view_path)
						panic("")
					}
					c.head = _template
				case "options":
					if c.options != nil {
						log.Error("You have two Options templates for %s, which is not allowed", view_path)
						panic("")
					}
					c.options = _template
				default:
					templateComponents[view_name+"."+file_name] = _template
					gotDefaultView = true
				}
				_template = nil
			}
		}
	}

	if !gotDefaultView {
		err := fmt.Errorf("default view not found for View %s in path %s | to fix this create a view with name default.html/php/gohtml or index.php/html/gohtml in the directory %s", view_path, view_path, view_path)
		panic(err)
	}
	return &c, nil
}
