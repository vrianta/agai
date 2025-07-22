package controller

import (
	"net/http"

	Session "github.com/vrianta/agai/v1/internal/session"
	Template "github.com/vrianta/agai/v1/template"
)

// Routes is a map of HTTP methods to their respective controllers
type (
	_Func      func(self *Context) *Template.Response // Map of HTTP methods to their respective handler functions
	_Templates struct {
		View    *Template.Struct // default template store
		Get     *Template.Struct // Template for GET requests
		POST    *Template.Struct // Template for POST requests
		DELETE  *Template.Struct // Template for DELETE requests
		PATCH   *Template.Struct // Template for PATCH requests
		PUT     *Template.Struct // Template for PUT requests
		HEAD    *Template.Struct // Template for HEAD requests
		OPTIONS *Template.Struct // Template for OPTIONS requests
	}

	Context struct {
		View      string     // name of the View have to mention it at the begining
		templates _Templates // storing pointer to a Template Struct store execute struct
		// HTTP methods with their respective handlers
		// Each method returns a view string and TemplateData
		// string is the template name to render and TemplateData is the data to pass to the template
		GET     _Func
		POST    _Func
		DELETE  _Func
		PATCH   _Func
		PUT     _Func
		HEAD    _Func
		OPTIONS _Func

		session *Session.Instance // Session object to handle user session

		// privte objects
		w http.ResponseWriter
		r *http.Request

		userInputs map[string]interface{}
	}
)
