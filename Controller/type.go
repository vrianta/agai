package Controller

import (
	"github.com/vrianta/Server/Session"
	"github.com/vrianta/Server/Template"
)

// Routes is a map of HTTP methods to their respective controllers
type (
	Struct struct {
		View string // name of the View have to mention it at the begining
		// HTTP methods with their respective handlers
		// Each method returns a view string and TemplateData
		// string is the template name to render and TemplateData is the data to pass to the template
		GET     func(*Session.Struct) *Template.Response
		POST    func(*Session.Struct) *Template.Response
		DELETE  func(*Session.Struct) *Template.Response
		PATCH   func(*Session.Struct) *Template.Response
		PUT     func(*Session.Struct) *Template.Response
		HEAD    func(*Session.Struct) *Template.Response
		OPTIONS func(*Session.Struct) *Template.Response
	}
)
