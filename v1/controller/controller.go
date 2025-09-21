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
	"net/http"

	"github.com/vrianta/agai/v1/config"
	"github.com/vrianta/agai/v1/internal/template"
	"github.com/vrianta/agai/v1/log"
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

	switch c.R.Method {
	case "GET":
		view := c.GET(c) // GET method of controller returns a view
		if view.asJson {
			// user want the response to be send as json
			c.W.Write(view.response.toJson())
			log.Debug("Template is nil for controller %s, no template to execute\n", view.name)
			return
		}
		__template, ok := template.GetTemplate(view.name)

		if !ok {
			log.Error("No Template Available as such name %s", view.name)
			panic("")
		}

		get_template := __template.GET()
		if !config.GetWebConfig().Build {
			get_template.Update()
		}
		if err := c.execute(get_template, view.response); err != nil {
			log.Error("Error rendering template: %T", err)
			panic(err)
		}
	case "POST":
		view := c.POST(c) // GET method of controller returns a view
		if view.asJson {
			// user want the response to be send as json
			c.W.Write(view.response.toJson())
			log.Debug("Template is nil for controller %s, no template to execute\n", view.name)
			return
		}
		__template, ok := template.GetTemplate(view.name)

		if !ok {
			log.Error("No Template Available as such name %s", view.name)
			panic("")
		}

		get_template := __template.POST()
		if !config.GetWebConfig().Build {
			get_template.Update()
		}
		if err := c.execute(get_template, view.response); err != nil {
			log.Error("Error rendering template: %T", err)
			panic(err)
		}
	case "DELETE":
		view := c.DELETE(c) // GET method of controller returns a view
		if view.asJson {
			// user want the response to be send as json
			c.W.Write(view.response.toJson())
			log.Debug("Template is nil for controller %s, no template to execute\n", view.name)
			return
		}
		__template, ok := template.GetTemplate(view.name)

		if !ok {
			log.Error("No Template Available as such name %s", view.name)
			panic("")
		}

		get_template := __template.DELETE()
		if !config.GetWebConfig().Build {
			get_template.Update()
		}
		if err := c.execute(get_template, view.response); err != nil {
			log.Error("Error rendering template: %T", err)
			panic(err)
		}
	case "PATCH":
		view := c.PATCH(c) // GET method of controller returns a view
		if view.asJson {
			// user want the response to be send as json
			c.W.Write(view.response.toJson())
			log.Debug("Template is nil for controller %s, no template to execute\n", view.name)
			return
		}
		__template, ok := template.GetTemplate(view.name)

		if !ok {
			log.Error("No Template Available as such name %s", view.name)
			panic("")
		}

		get_template := __template.PATCH()
		if !config.GetWebConfig().Build {
			get_template.Update()
		}
		if err := c.execute(get_template, view.response); err != nil {
			log.Error("Error rendering template: %T", err)
			panic(err)
		}
	case "PUT":
		view := c.PUT(c) // GET method of controller returns a view
		if view.asJson {
			// user want the response to be send as json
			c.W.Write(view.response.toJson())
			log.Debug("Template is nil for controller %s, no template to execute\n", view.name)
			return
		}
		__template, ok := template.GetTemplate(view.name)

		if !ok {
			log.Error("No Template Available as such name %s", view.name)
			panic("")
		}

		get_template := __template.PUT()
		if !config.GetWebConfig().Build {
			get_template.Update()
		}
		if err := c.execute(get_template, view.response); err != nil {
			log.Error("Error rendering template: %T", err)
			panic(err)
		}
	case "HEAD":
		view := c.HEAD(c) // GET method of controller returns a view
		if view.asJson {
			// user want the response to be send as json
			c.W.Write(view.response.toJson())
			log.Debug("Template is nil for controller %s, no template to execute\n", view.name)
			return
		}
		__template, ok := template.GetTemplate(view.name)

		if !ok {
			log.Error("No Template Available as such name %s", view.name)
			panic("")
		}

		get_template := __template.HEAD()
		if !config.GetWebConfig().Build {
			get_template.Update()
		}
		if err := c.execute(get_template, view.response); err != nil {
			log.Error("Error rendering template: %T", err)
			panic(err)
		}
	case "OPTIONS":
		view := c.OPTIONS(c) // GET method of controller returns a view
		if view.asJson {
			// user want the response to be send as json
			c.W.Write(view.response.toJson())
			log.Debug("Template is nil for controller %s, no template to execute\n", view.name)
			return
		}
		__template, ok := template.GetTemplate(view.name)

		if !ok {
			log.Error("No Template Available as such name %s", view.name)
			panic("")
		}

		get_template := __template.OPTIONS()
		if !config.GetWebConfig().Build {
			get_template.Update()
		}
		if err := c.execute(get_template, view.response); err != nil {
			log.Error("Error rendering template: %T", err)
			panic(err)
		}
	default:
		view := c.GET(c) // GET method of controller returns a view
		if view.asJson {
			// user want the response to be send as json
			c.W.Write(view.response.toJson())
			log.Debug("Template is nil for controller %s, no template to execute\n", view.name)
			return
		}
		__template, ok := template.GetTemplate(view.name)

		if !ok {
			log.Error("No Template Available as such name %s", view.name)
			panic("")
		}

		get_template := __template.INDEX()
		if !config.GetWebConfig().Build {
			get_template.Update()
		}
		if err := c.execute(get_template, view.response); err != nil {
			log.Error("Error rendering template: %T", err)
			panic(err)
		}
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
	c.W = w
	c.R = r

	c.runRequest()

}

func (c *Context) GET(self *Context) view {
	response := Response{
		"Massage": "Welcome to Defalut Get Page",
	}
	return response.AsJson()
}

func (c *Context) POST(self *Context) view {
	return c.GET(self)
}

func (c *Context) PUT(self *Context) view {
	return c.GET(self)
}

func (c *Context) PATCH(self *Context) view {
	return c.GET(self)
}

func (c *Context) DELETE(self *Context) view {
	return c.GET(self)
}

func (c *Context) HEAD(self *Context) view {
	return c.GET(self)
}

func (c *Context) OPTIONS(self *Context) view {
	return c.GET(self)
}
