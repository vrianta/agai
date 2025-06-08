package Controller

import (
	"fmt"
	"net/http"
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
	if session != nil {
		session.Update(c.w, c.r)
	}
	switch c.r.Method {
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
	if _template, err := Template.New(c.View); err != nil {
		return err
	} else {
		c.template = _template
		return nil
	}
}

func (c *Struct) Execute(__response *Template.Response) error {
	if Config.Build {
		return c.template.Execute(c.w, __response)
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
