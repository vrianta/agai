package router

import (
	"github.com/vrianta/agai/v1/internal/requestHandler"
)

type context struct {
	root string
}

func New(root string) *context {
	return &context{
		root: root,
	}
}

func (c *context) AddRoute(route string, controller requestHandler.ControllerInterface) *context {
	requestHandler.RouteTable[c.root+"/"+route] = controller

	return c
}

func CreateRoute[T requestHandler.ControllerInterface](route string, controller T) {

	requestHandler.RouteTable[route] = controller
}
