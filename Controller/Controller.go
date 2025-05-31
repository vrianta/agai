package Controller

import (
	"fmt"

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
func (c *Struct) CallMethod(method string, session *Session.Struct) *Template.Response {
	switch method {
	case "GET":
		return c.isMethodNull(c.GET, session)
	case "POST":
		return c.isMethodNull(c.POST, session)
	case "DELETE":
		return c.isMethodNull(c.DELETE, session)
	case "PATCH":
		return c.isMethodNull(c.PATCH, session)
	case "PUT":
		return c.isMethodNull(c.PUT, session)
	case "HEAD":
		return c.isMethodNull(c.HEAD, session)
	case "OPTIONS":
		return c.isMethodNull(c.OPTIONS, session)
	default:
		return &Template.Response{"error": "Method not allowed"}
	}
}

/*
 * This Methid is to check if the Method passing is Defined or not if nill will return Error else print the value
 */
func (c *Struct) isMethodNull(method func(*Session.Struct) *Template.Response, session *Session.Struct) *Template.Response {
	if method != nil {
		return method(session)
	}

	return &Template.Response{"error": "Current Method is not allowed"}
}

func (c *Struct) Validate() {
	if c.View == "" {
		panic(fmt.Errorf("view is not defined for the controller %T", c))
	}

}
