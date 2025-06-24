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
This file will store the local method for Controller which will be by default and Entirely local
*/

/*
 * Controller Function Call
 * This function will be responsible for handling the Method calling of te Controller
 * Example: if the Method is GET then it will call Get Method if The Method is POST then it will call Post Method
 * This will be used in the routingHandler to call the correct method of the controller
 */
func (c *Struct) RunRequest(session *Session.Struct) {
	c.assignSession(session) // Assign the session to the controller
	if session != nil {
		session.Update(c.w, c.r)
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
		// fmt.write(c.w, "Method Not Allowed")
	}
}

/*
 * This Methid is to check if the Method passing is Defined or not if nill will return Error else print the value
 */
func (c *Struct) isMethodNull(method _Func) *Template.Response {
	if method != nil {
		return method(c)
	}

	return &Template.Response{"error": "Current Method is not allowed"}
}

// Function to Assign the Session in the Controller
func (c *Struct) assignSession(session *Session.Struct) {
	c.session = session
}

func (c *Struct) Validate() {
	if c.View == "" {
		panic(fmt.Errorf("view is not defined for the controller %T", c))
	}
}

// Function to return the Session because do not want to expose the session varaible directly
func (c *Struct) GetSession() *Session.Struct {
	return c.session
}

// Function which will check the View and it's extension and determine the type of view and accordingly will call the
// appropiate render function
func (c *Struct) RenderView(view string, data *Template.Response) error {
	if view == "" {
		return nil // No view to render
	}

	// get the extension of the view
	extension := strings.Split(view, ".")

	switch extension {
	// case "html", "htm", "gohtml":
	}

	// if err := c.Session.RenderEngine.RenderTemplate(view, data); err != nil {
	// 	return fmt.Errorf("error rendering view %s: %w", view, err)
	// }
	return nil
}

func (c *Struct) RegisterTemplate() error {

	// Get all the Files in the View Folder with certain Name suffix
	// Views should be inside a folder named as the view name and the views names with be with certain specific name types
	// example for default use default.html or default.php // same file with two extensions are not allowed
	// for Get Method use get.html or get.php

	view_path := "./" + Config.ViewFolder + "/" + c.View // view path of the view package of the controller
	files, err := os.ReadDir(view_path)
	if err != nil {
		err := fmt.Errorf("error reading directory: %s", err.Error())
		panic(err)
	}

	var gotDefaultView = false // to check if we got the default view or not
	for _, entry := range files {
		if !entry.IsDir() {
			full_file_name := entry.Name()
			var file_type = strings.TrimPrefix(filepath.Ext(full_file_name), ".") // type of the file
			file_name := full_file_name[:len(full_file_name)-len(file_type)-1]

			// Template is our Custom Template Package this is not the go one
			// Get name without extension
			// fmt.Println("File Name:", file_name, "File Type:", file_type)
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

func (c *Struct) ExecuteTemplate(__template *Template.Struct, __response *Template.Response) error {
	if !Config.Build {
		return __template.Update()
	}

	if __template == nil {
		if err := c.templates.View.Execute(c.w, __response); err != nil {
			Log.WriteLog("Error rendering template: " + err.Error())
			panic(err)
		}
	} else {
		if err := __template.Execute(c.w, __response); err != nil {
			Log.WriteLog("Error rendering template: " + err.Error())
			panic(err)
		}
	}

	// return __template.Execute(c.w, __response)
	return nil
}

// a Copy Function to create a new controller Instance by copying the data
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

		session: c.session,

		// userInputs: make(map[string]interface{}, 20),
	}
}

/*
 * To Initialise the controller with the writer and reader objects
 */
func (c *Struct) InitWR(w http.ResponseWriter, r *http.Request) {
	c.w = w
	c.r = r
}

/*
 * To Initialise the controller with the Session
 */
func (c *Struct) InitSession(__s *Session.Struct) {
	c.session = __s
}
