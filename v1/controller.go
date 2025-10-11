package agai

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

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/vrianta/agai/v1/internal/session"
	"github.com/vrianta/agai/v1/log"
	"github.com/vrianta/agai/v1/utils"
)

// Routes is a map of HTTP methods to their respective controllers
type (
	Controller struct {
		session *session.Instance // Session object to handle user session

		// privte objects
		W http.ResponseWriter
		R *http.Request

		userInputs map[string]any
	}

	View = func() *view
)

func (c *Controller) init(w http.ResponseWriter, r *http.Request) {
	c.W = w
	c.R = r
}

/*
This file contains local methods for the Controller struct, providing default and internal logic.
*/
func (c *Controller) GET() View {
	response := Response{
		"Massage": "Welcome to Defalut Get Page",
	}
	return c.ResponseAsJson(response)
}

func (c *Controller) POST() View {
	return c.GET()
}

func (c *Controller) PUT() View {
	return c.GET()
}

func (c *Controller) PATCH() View {
	return c.GET()
}

func (c *Controller) DELETE() View {
	return c.GET()
}

func (c *Controller) HEAD() View {
	return c.GET()
}

func (c *Controller) OPTIONS() View {
	return c.GET()
}

/*
 * Store Data in the Session
 */
func (_c *Controller) StoreData(index string, _d any) {
	_c.session.Store(index, _d)
}

/*
 * Get Data From Session Store
 */
func (_c *Controller) GetStoredData(index string) (any, bool) {
	data, ok := _c.session.Data[index]
	return data, ok
}

// Return all Inputs at once
func (_c *Controller) GetInputs() *map[string]any {
	if _c.userInputs == nil {
		_c.parseRequest()
	}
	return &_c.userInputs
}

// Return specific input
// if present then value else nil
func (_c *Controller) GetInput(key string) (any, error) {
	if _c.userInputs == nil {
		_c.parseRequest()
	}
	if v, ok := _c.userInputs[key]; ok {
		return v, nil
	}
	return nil, errors.New("NOINPUT")
}

/*
 * This File is to handle User Inputs
 */
func (_c *Controller) parseRequest() {

	_c.userInputs = make(map[string]any, 30)

	// Log handling of queryBuilder parameters for non-POST methods
	if _c.R.Method == http.MethodGet {
		for key, values := range _c.R.URL.Query() {
			_c.processqueryBuilderParams(key, values)
		}
	}

	contentType := _c.R.Header.Get("Content-Type")
	switch contentType {
	case "application/json":
		if p, err := io.ReadAll(_c.R.Body); err != nil {
			log.Error("Failed to Read the Joson Data, %v\n", err)
		} else {
			if er := json.Unmarshal(p, &_c.userInputs); er != nil {
				log.Error("Failed to Render the Json Data, %v\n", er)
			}
		}
	case "application/x-www-form-urlencoded":
		// Handle form data (application/x-www-form-urlencoded)
		if err := _c.R.ParseForm(); err != nil {
			log.WriteLogf("Error parsing form data | Error - %s\n", err.Error())
		} else {
			for key, values := range _c.R.PostForm {
				_c.processPostParams(key, values)
			}
		}

	case "multipart/form-data":
		// Handle multipart form data (file upload)
		// Note: This case is handled separately below
		if err := _c.R.ParseMultipartForm(10 << 20); err != nil { // 10 MB
			log.WriteLogf("Error parsing multipart form data | Error - %s\n", err.Error())
			return
		}
		for key, values := range _c.R.PostForm {
			_c.processPostParams(key, values)
		}

	default:
		log.Error("Content-Type %s not supported yet raise a issue in github to get it implimented", contentType)
	}

}

// handlequeryBuilderParams processes parameters found in the URL queryBuilder
func (_c *Controller) processqueryBuilderParams(key string, values []string) {
	var err error
	// Check for multiple values

	if len(values) > 1 {
		if _c.userInputs[key], err = utils.JsonToString(values); err != nil {
			// http.Error(sh.W, "Failed to convert data to JSON", http.StatusMethodNotAllowed)
			return
		}
	} else {
		_c.userInputs[key] = values[0] // Store single value as a string
	}
}

// handlePostParams processes parameters found in the POST data
func (_c *Controller) processPostParams(key string, values []string) {
	var err error
	if len(values) > 1 {
		if _c.userInputs[key], err = utils.JsonToString(values); err != nil {
			// http.Error(sh.W, "Failed to convert data to JSON", http.StatusMethodNotAllowed)
			return
		}
	} else {
		_c.userInputs[key] = values[0] // Store single value as a string
	}
}
