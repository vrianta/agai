package requestHandler

import (
	"net/http"

	"github.com/vrianta/agai/v1/controller"
)

// // Handler processes incoming HTTP requests and manages user sessions.
// // It checks if the user has an existing session and handles session creation or validation.
// // Based on the session and route, it invokes the appropriate controller method.
// // Parameters:
// // - w: The HTTP response writer.
// // - r: The HTTP request.
type (
	RouteDestination interface { // Resembeles Controller Package
		GET(self *controller.Context)
		POST(self *controller.Context)
		PUT(self *controller.Context)
		DELETE(self *controller.Context)
		PATCH(self *controller.Context)
		HEAD(self *controller.Context)
		OPTIONS(self *controller.Context)
		RunRequest()
	}

	routes map[string]func(w http.ResponseWriter, r *http.Request) RouteDestination
)

var RouteTable routes

func Handler(w http.ResponseWriter, r *http.Request) {

	if _c, found := RouteTable[r.URL.Path]; found {
		tempController := _c(w, r)
		tempController.RunRequest()
	} else {
		http.Error(w, "404 Error : Route not found ", http.StatusNotFound)
		return
	}
}
