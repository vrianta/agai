package agai

import "net/http"

// it will keep a record if the user has registered / root path or not
var RootRegistered bool = false

func CreateRoute[T any, PT interface {
	*T
	controllerInterface
}](route string) {

	if route == "/" {
		RootRegistered = true
		http.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/" {
				http.Redirect(w, r, "/404", int(HttpStatus.SeeOther))
				return
			}
			var tempController PT = new(T)
			tempController.init(w, r)
			runRequest(w, r, tempController)
		})
		return
	}

	http.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		var tempController PT = new(T)
		tempController.init(w, r)
		runRequest(w, r, tempController)
	})
}
