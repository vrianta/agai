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
)

func (c *Context) Init(w http.ResponseWriter, r *http.Request) {
	c.w = w
	c.r = r
}

/*
This file contains local methods for the Controller struct, providing default and internal logic.
*/

func (c *Context) GET() View {
	response := Response{
		"Massage": "Welcome to Defalut Get Page",
	}
	return response.AsJson()
}

func (c *Context) POST() View {
	response := Response{
		"Massage": "Welcome to Defalut Get Page",
	}
	return response.AsJson()
}

func (c *Context) PUT() View {
	response := Response{
		"Massage": "Welcome to Defalut Get Page",
	}
	return response.AsJson()
}

func (c *Context) PATCH() View {
	response := Response{
		"Massage": "Welcome to Defalut Get Page",
	}
	return response.AsJson()
}

func (c *Context) DELETE() View {
	response := Response{
		"Massage": "Welcome to Defalut Get Page",
	}
	return response.AsJson()
}

func (c *Context) HEAD() View {
	response := Response{
		"Massage": "Welcome to Defalut Get Page",
	}
	return response.AsJson()
}

func (c *Context) OPTIONS() View {
	response := Response{
		"Massage": "Welcome to Defalut Get Page",
	}
	return response.AsJson()
}
