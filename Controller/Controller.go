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
func (c *Struct) CallMethod() *Template.Response {
	switch c.Session.R.Method {
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
func (c *Struct) AssignSession(session *Session.Struct) {
	c.Session = session
}

func (c *Struct) Validate() {
	if c.View == "" {
		panic(fmt.Errorf("view is not defined for the controller %T", c))
	}
}
