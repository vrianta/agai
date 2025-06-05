package Controller

import (
	"github.com/vrianta/Server/Session"
	"github.com/vrianta/Server/Template"
)

// Routes is a map of HTTP methods to their respective controllers
type (
	_Func func(self *Struct) *Template.Response // Map of HTTP methods to their respective handler functions

	Struct struct {
		View     string           // name of the View have to mention it at the begining
		template *Template.Struct // storing pointer to a Template Struct store execute struct
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

		Session *Session.Struct // Session object to handle user session
	}
)
