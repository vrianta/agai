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
    - Initializing request/response/session objects (InitWR, initSession)
    - Running the correct handler for an HTTP request (runRequest)
    - Copying controller instances (Copy)
    - Accessing the session safely (GetSession)

Usage:
------
- Define a controller struct embedding Controller.Instance.
- Set the View field and handler functions for each HTTP method.
- Register the controller in your router.
- Use the provided methods to manage templates and session state.

See Also:
---------
- Session package: for session management and data storage.
- Template package: for template parsing and rendering.
- Config package: for server and view configuration.

*/

package controller

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	Config "github.com/vrianta/agai/v1/config"
	Template "github.com/vrianta/agai/v1/template"
)

/*
This file contains local methods for the Controller struct, providing default and internal logic.
*/

/*
runRequest dispatches the HTTP request to the appropriate handler method (GET, POST, etc.)
and renders the corresponding template. It also assigns and updates the session for the request.

Parameters:
- session: pointer to the current Session.Instance for the request.
*/
func (c *Context) runRequest() {

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
		c.ExecuteTemplate(c.templates.View, &Template.NoResponse)
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
func (c *Context) isMethodNull(method _Func) *Template.Response {
	if method != nil {
		return method(c)
	}
	return &Template.Response{"error": "Current Method is not allowed"}
}

/*
Validate checks if the controller's View field is set.
Panics if the View is not defined, ensuring every controller has an associated view.
*/
func (c *Context) Validate() {
	if c.View == "" {
		panic(fmt.Errorf("view is not defined for the controller %T", c))
	}
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
func (c *Context) RenderView(view string, data *Template.Response) error {
	if view == "" {
		return nil // No view to render
	}
	extension := strings.Split(view, ".")
	// Extend this switch to support more view types if needed
	switch extension {
	// case "html", "htm", "gohtml":
	}
	return nil
}

/*
RegisterTemplate scans the controller's view directory and registers templates for each HTTP method.
It expects files named default.html/php/gohtml, get.html/php, post.html/php, etc.
Panics if no default view is found.

Returns:
- error: if reading the directory or registering a template fails.
*/
func (c *Context) RegisterTemplate() error {
	view_path := "./" + Config.GetWebConfig().ViewFolder + "/" + c.View // Path to the controller's view folder
	// fmt.Printf("Registering templates for controller: %T, view path: %s\n", c, view_path)
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
				// fmt.Printf("  Found template: %s (type: %s) for controller: %T and file_name:%s Path:%s\n", full_file_name, file_type, c, file_name, view_path)
				switch file_name {
				case "default", "index":
					// fmt.Printf("  Registering default view template for controller: %T\n", c)
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
Copy creates a new Instance of the controller with the same configuration and handlers.
Useful for creating per-request controller instances.

Returns:
- *Instance: pointer to the copied controller struct.
*/
func (c *Context) Copy() *Context {
	return &Context{
		View:      c.View,
		templates: c.templates,
		GET:       c.GET,
		POST:      c.POST,
		DELETE:    c.DELETE,
		PATCH:     c.PATCH,
		PUT:       c.PUT,
		HEAD:      c.HEAD,
		OPTIONS:   c.OPTIONS,
	}
}

/*
InitWR initializes the controller with the HTTP response writer and request.
Call this before handling a request.

Parameters:
- w: http.ResponseWriter for the response.
- r: *http.Request for the incoming request.
*/
func (c *Context) Init(w http.ResponseWriter, r *http.Request) {
	c.w = w
	c.r = r

	c.runRequest()

}
