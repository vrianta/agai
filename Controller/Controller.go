/*
Package Controller

This package defines the core Controller struct and methods for the Go Server Framework.
Controllers are responsible for handling HTTP requests, managing session state, and rendering views.

Key Concepts:
-------------
- Each Controller is a struct that maps HTTP methods (GET, POST, etc.) to handler functions.
- The Controller manages its own session, request, and response writer.
- Views are loaded from the configured Views folder, and templates are registered per HTTP method.
- Templates support PHP-like syntax and are rendered using the custom Template engine.
- Session data is accessed via the controller's session field (see Session package).
- The Controller provides helper methods for:
    - Validating configuration (Validate)
    - Registering and executing templates (RegisterTemplate, ExecuteTemplate)
    - Initializing request/response/session objects (InitWR, InitSession)
    - Running the correct handler for an HTTP request (RunRequest)
    - Copying controller instances (Copy)
    - Accessing the session safely (GetSession)

Usage:
------
- Define a controller struct embedding Controller.Struct.
- Set the View field and handler functions for each HTTP method.
- Register the controller in your router.
- Use the provided methods to manage templates and session state.

See Also:
---------
- Session package: for session management and data storage.
- Template package: for template parsing and rendering.
- Config package: for server and view configuration.

*/

package Controller

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/vrianta/Server/Config"
	"github.com/vrianta/Server/Log"
	"github.com/vrianta/Server/Session"
	"github.com/vrianta/Server/Template"
)

/*
This file contains local methods for the Controller struct, providing default and internal logic.
*/

/*
RunRequest dispatches the HTTP request to the appropriate handler method (GET, POST, etc.)
and renders the corresponding template. It also assigns and updates the session for the request.

Parameters:
- session: pointer to the current Session.Struct for the request.
*/
func (c *Struct) RunRequest(session *Session.Struct) {
	c.assignSession(session) // Assign the session to the controller
	if session != nil {
		session.Update(c.w, c.r) // Update session with current writer and request
	}
	switch c.r.Method {
	case "GET":
		reponse := c.isMethodNull(c.GET)
		c.ExecuteTemplate(c.templates.Get, reponse)
	case "POST":
		reponse := c.isMethodNull(c.POST)
		c.ExecuteTemplate(c.templates.POST, reponse)
	case "DELETE":
		reponse := c.isMethodNull(c.DELETE)
		c.ExecuteTemplate(c.templates.DELETE, reponse)
	case "PATCH":
		reponse := c.isMethodNull(c.PATCH)
		c.ExecuteTemplate(c.templates.PATCH, reponse)
	case "PUT":
		reponse := c.isMethodNull(c.PUT)
		c.ExecuteTemplate(c.templates.PUT, reponse)
	case "HEAD":
		reponse := c.isMethodNull(c.HEAD)
		c.ExecuteTemplate(c.templates.HEAD, reponse)
	case "OPTIONS":
		reponse := c.isMethodNull(c.OPTIONS)
		c.ExecuteTemplate(c.templates.OPTIONS, reponse)
	default:
		// Method not allowed; optionally handle here
	}
}

/*
isMethodNull checks if the provided handler function is nil.
If not nil, it calls the handler and returns its response.
If nil, returns a default error response.

Parameters:
- method: the handler function for the HTTP method.

Returns:
- *Template.Response: the response to render.
*/
func (c *Struct) isMethodNull(method _Func) *Template.Response {
	if method != nil {
		return method(c)
	}
	return &Template.Response{"error": "Current Method is not allowed"}
}

/*
assignSession assigns the given session to the controller instance.

Parameters:
- session: pointer to Session.Struct to assign.
*/
func (c *Struct) assignSession(session *Session.Struct) {
	c.session = session
}

/*
Validate checks if the controller's View field is set.
Panics if the View is not defined, ensuring every controller has an associated view.
*/
func (c *Struct) Validate() {
	if c.View == "" {
		panic(fmt.Errorf("view is not defined for the controller %T", c))
	}
}

/*
GetSession safely returns the controller's session pointer.
Use this instead of accessing the session field directly.
*/
func (c *Struct) GetSession() *Session.Struct {
	return c.session
}

/*
RenderView determines the type of view based on its extension and calls the appropriate render function.
Currently a stub; extend this to support multiple template engines if needed.

Parameters:
- view: the view file name.
- data: pointer to Template.Response containing data for the template.

Returns:
- error: if rendering fails.
*/
func (c *Struct) RenderView(view string, data *Template.Response) error {
	if view == "" {
		return nil // No view to render
	}
	extension := strings.Split(view, ".")
	// Extend this switch to support more view types if needed
	switch extension {
	// case "html", "htm", "gohtml":
	}
	// Example for future use:
	// if err := c.Session.RenderEngine.RenderTemplate(view, data); err != nil {
	// 	return fmt.Errorf("error rendering view %s: %w", view, err)
	// }
	return nil
}

/*
RegisterTemplate scans the controller's view directory and registers templates for each HTTP method.
It expects files named default.html/php/gohtml, get.html/php, post.html/php, etc.
Panics if no default view is found.

Returns:
- error: if reading the directory or registering a template fails.
*/
func (c *Struct) RegisterTemplate() error {
	view_path := "./" + Config.GetWebConfig().ViewFolder + "/" + c.View // Path to the controller's view folder
	files, err := os.ReadDir(view_path)
	if err != nil {
		err := fmt.Errorf("error reading directory: %s", err.Error())
		panic(err)
	}

	var gotDefaultView = false // Track if a default view is found
	for _, entry := range files {
		if !entry.IsDir() {
			full_file_name := entry.Name()
			var file_type = strings.TrimPrefix(filepath.Ext(full_file_name), ".") // File extension/type
			file_name := full_file_name[:len(full_file_name)-len(file_type)-1]    // Name without extension

			// Register the template using the custom Template package
			if _template, err := Template.New(view_path, full_file_name, file_type); err != nil {
				return err
			} else {
				switch file_name {
				case "default", "index":
					c.templates.View = _template
					gotDefaultView = true
				case "get":
					c.templates.Get = _template
				case "post":
					c.templates.POST = _template
				case "delete":
					c.templates.DELETE = _template
				case "patch":
					c.templates.PATCH = _template
				case "put":
					c.templates.PUT = _template
				case "head":
					c.templates.HEAD = _template
				case "options":
					c.templates.OPTIONS = _template
				default:
					// Ignore unknown files
				}
			}
		}
	}

	if !gotDefaultView {
		err := fmt.Errorf("default view not found for controller %s in path %s | to fix this create a view with name default.html/php/gohtml or index.php/html/gohtml in the directory %s", c.View, view_path, view_path)
		panic(err)
	}
	return nil
}

/*
ExecuteTemplate renders the given template with the provided response data.
If not in build mode, updates the template before rendering.
Logs and panics on rendering errors.

Parameters:
- __template: pointer to the Template.Struct to render.
- __response: pointer to Template.Response containing data for the template.

Returns:
- error: if updating the template fails (in dev mode).
*/
func (c *Struct) ExecuteTemplate(__template *Template.Struct, __response *Template.Response) error {
	if __template == nil {
		__template = c.templates.View
	}
	if !Config.GetWebConfig().Build {
		__template.Update()
		if err := __template.Execute(c.w, __response); err != nil {
			Log.WriteLog("Error rendering template: " + err.Error())
			panic(err)
		} // Update the template before rendering
		return nil
	}
	if c.View == "" {
		return nil // No view to render, return nil
	}

	// if __template == nil {
	// 	if err := c.templates.View.Execute(c.w, __response); err != nil {
	// 		Log.WriteLog("Error rendering template: " + err.Error())
	// 		panic(err)
	// 	}
	// } else {
	if err := __template.Execute(c.w, __response); err != nil {
		Log.WriteLog("Error rendering template: " + err.Error())
		return err
	}
	// }
	return nil
}

/*
Copy creates a new instance of the controller with the same configuration and handlers.
Useful for creating per-request controller instances.

Returns:
- *Struct: pointer to the copied controller struct.
*/
func (c *Struct) Copy() *Struct {
	return &Struct{
		View:      c.View,
		templates: c.templates,
		GET:       c.GET,
		POST:      c.POST,
		DELETE:    c.DELETE,
		PATCH:     c.PATCH,
		PUT:       c.PUT,
		HEAD:      c.HEAD,
		OPTIONS:   c.OPTIONS,
		session:   c.session,
		// userInputs: make(map[string]interface{}, 20),
	}
}

/*
InitWR initializes the controller with the HTTP response writer and request.
Call this before handling a request.

Parameters:
- w: http.ResponseWriter for the response.
- r: *http.Request for the incoming request.
*/
func (c *Struct) InitWR(w http.ResponseWriter, r *http.Request) {
	c.w = w
	c.r = r
}

/*
InitSession assigns the given session to the controller.
Call this to set up session state for the request.

Parameters:
- __s: pointer to Session.Struct to assign.
*/
func (c *Struct) InitSession(__s *Session.Struct) {
	c.session = __s
}
