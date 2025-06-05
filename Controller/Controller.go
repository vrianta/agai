package Controller

import (
	"fmt"
	"strings"

	"github.com/vrianta/Server/Config"
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
func (c *Struct) CallMethod(session *Session.Struct) *Template.Response {
	c.assignSession(session) // Assign the session to the controller
	switch session.R.Method {
	case "GET":
		return c.isMethodNull(c.GET)
	case "POST":
		return c.isMethodNull(c.POST)
	case "DELETE":
		return c.isMethodNull(c.DELETE)
	case "PATCH":
		return c.isMethodNull(c.PATCH)
	case "PUT":
		return c.isMethodNull(c.PUT)
	case "HEAD":
		return c.isMethodNull(c.HEAD)
	case "OPTIONS":
		return c.isMethodNull(c.OPTIONS)
	default:
		return &Template.Response{"error": "Method not allowed"}
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
	c.Session = session
}

func (c *Struct) Validate() {
	if c.View == "" {
		panic(fmt.Errorf("view is not defined for the controller %T", c))
	}
}

// Function to return the Session because do not want to expose the session varaible directly
func (c *Struct) GetSession() *Session.Struct {
	return c.Session
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
	if _template, err := Template.New(c.View); err != nil {
		return err
	} else {
		c.template = _template
		return nil
	}
}

func (c *Struct) Execute(__response *Template.Response) error {
	if Config.Build {
		return c.template.Execute(c.Session.W, __response)
	}

	return c.template.Update()

}

// a Copy Function to create a new controller Instance by copying the data
func (c *Struct) Copy() *Struct {
	return &Struct{
		View:     c.View,
		template: c.template,
		GET:      c.GET,
		POST:     c.POST,
		DELETE:   c.DELETE,
		PATCH:    c.PATCH,
		PUT:      c.PUT,
		HEAD:     c.HEAD,
		OPTIONS:  c.OPTIONS,

		Session: c.Session,
	}
}
