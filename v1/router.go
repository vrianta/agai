package agai

import "net/http"

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
	controllerInterface
}](route string) {

	http.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		var tempController PT = new(T)
		tempController.init(w, r)
		runRequest(w, r, tempController)
	})
}
