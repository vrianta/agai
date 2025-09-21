package requesthandler

import (
	"net/http"

	"github.com/vrianta/agai/v1/controller"
)

var routeTable routes // map[string]*Controller.Struct

type (
	// resembles the controller interface
	routeDestination interface {
		GET(self *controller.Context)
		POST(self *controller.Context)
		PUT(self *controller.Context)
		DELETE(self *controller.Context)
		PATCH(self *controller.Context)
		HEAD(self *controller.Context)
		OPTIONS(self *controller.Context)
		Init(w http.ResponseWriter, r *http.Request)
		runRequest()
	}

	routes map[string]func(w http.ResponseWriter, r *http.Request) routeDestination // Type for routes, mapping URL paths to Controller structs

)

func CreateRoute[T routeDestination](route string) {

	routeTable[route] = func(w http.ResponseWriter, r *http.Request) routeDestination {
		var t T
		return t
	}
}

//
// // Handler processes incoming HTTP requests and manages user sessions.
// // It checks if the user has an existing session and handles session creation or validation.
// // Based on the session and route, it invokes the appropriate controller method.
// // Parameters:
// // - w: The HTTP response writer.
// // - r: The HTTP request.
// func Handler(w http.ResponseWriter, r *http.Request) {

// 	if _controller, found := routeTable[r.URL.Path]; found {
// 		tempController := _controller{
// 			W: w
// 		}
// 		tempController.Init(w, r)
// 	} else {
// 		http.Error(w, "404 Error : Route not found ", http.StatusNotFound)
// 		return
// 	}
// }
