package controller

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/vrianta/agai/v1/log"
	"github.com/vrianta/agai/v1/utils"
)

/*
 * This File is to handle User Inputs
 */

func (_c *Context) ParseRequest() {

	_c.userInputs = make(map[string]interface{}, 30)

	// Log handling of queryBuilder parameters for non-POST methods
	for key, values := range _c.R.URL.Query() {
		_c.processqueryBuilderParams(key, values)
	}

	if _c.R.Method == http.MethodPost {
		contentType := _c.R.Header.Get("Content-Type")
		switch contentType {
		case "application/json":
			// log.Debug("Got Json Body")
			// b :=

			if p, err := io.ReadAll(_c.R.Body); err != nil {
				log.Error("Failed to Read the Joson Data, %v", err)
			} else {
				if er := json.Unmarshal(p, &_c.userInputs); er != nil {
					log.Error("Failed to Render the Json Data, %v", er)
				}
			}
			// fmt.Println(string(p))
			// fmt.Println(_c.userInputs)

		case "application/x-www-form-urlencoded":
			// Handle form data (application/x-www-form-urlencoded)
			err := _c.R.ParseForm()
			if err != nil {
				log.WriteLogf("Error parsing form data | Error - %s", err.Error())
				return
			}
			for key, values := range _c.R.PostForm {
				_c.processPostParams(key, values)
			}

		case "multipart/form-data":
			// Handle multipart form data (file upload)
			// Note: This case is handled separately below
			if err := _c.R.ParseMultipartForm(10 << 20); err != nil { // 10 MB
				log.WriteLogf("Error parsing multipart form data | Error - %s", err.Error())
				return
			}
			for key, values := range _c.R.PostForm {
				_c.processPostParams(key, values)
			}

		default:
			break
		}

		// Log handling of queryBuilder parameters for non-POST methods
		// _c.userInputs = _c.R.PostForm

	}

}

// handlequeryBuilderParams processes parameters found in the URL queryBuilder
func (_c *Context) processqueryBuilderParams(key string, values []string) {
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
func (_c *Context) processPostParams(key string, values []string) {
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

// Return all Inputs at once
func (_c *Context) GetInputs() *map[string]interface{} {
	if _c.userInputs == nil {
		_c.ParseRequest()
	}
	return &_c.userInputs
}

// Return specific input
// if present then value else nil
func (_c *Context) GetInput(key string) interface{} {
	if _c.userInputs == nil {
		_c.ParseRequest()
	}
	if v, ok := _c.userInputs[key]; ok {
		return v
	}
	return nil
}
