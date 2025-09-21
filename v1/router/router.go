package router

import (
	"net/http"

	"github.com/vrianta/agai/v1/internal/requestHandler"
)

// /*
//  * Create a New Router Object with Default route group example / is the default route for this or /api or /v1 etc
//  */
// func New(root string) {
// 	defaultRoute = root
// }

func CreateRoute[T requestHandler.RouteDestination](route string) {

	requestHandler.RouteTable[route] = func(w http.ResponseWriter, r *http.Request) requestHandler.RouteDestination {
		var t T
		return t
	}
}

// A Function to Create and Return
