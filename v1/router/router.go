package router

import (
	"github.com/vrianta/agai/v1/internal/requestHandler"
)

type context struct {
	root string
}

// func New(root string) *context {
// 	return &context{
// 		root: root,
// 	}
// }

// func (c *context) AddRoute[T any, PT interface {
// 	*T
// 	requestHandler.ControllerInterface
// }](route string, controller requestHandler.ControllerInterface) *context {
// 	requestHandler.RouteTable[c.root+"/"+route] = func() requestHandler.ControllerInterface {
// 		var c PT = new(T)
// 		return c
// 	}

// 	return c
// }

func CreateRoute[T any, PT interface {
	*T
	requestHandler.ControllerInterface
}](route string) {

	requestHandler.RouteTable[route] = func() requestHandler.ControllerInterface {
		var c PT = new(T)
		return c
	}
	// requestHandler.RouteTable[route] = controller
}
