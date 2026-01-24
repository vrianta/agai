package agai

import (
	"net/http"
	"strings"
)

// it will keep a record if the user has registered / root path or not
var RootRegistered bool = false
var rootPath string = ""

/*
 * If you want to change the initial path for next registered paths then pass paths
 * example you want next all paths to be /example/then-other-paths then pass example in the function
 * GroupPathWith("example")
 * As is you can pass multiple paths to be initial path
 * GroupPathWith("example", "path1") -> /example/path1
 * CreateRoute[Controller]("login") -> /example/path1/login/
 */
func GroupNextPathsWith(route ...string) {
	rootPath = "/" + strings.Join(route, "/")
}

/*
 * CreateRoute("login", "register")
 * it will create two paths
 * 1. /login/register/
 * 2. /login/register
 */
func CreateRoute[T any, PT interface {
	*T
	controllerInterface
}](route ...string) {

	if route[0] == "/" && (rootPath == "" || rootPath == "/") {
		RootRegistered = true
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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

	if (route[0] == "" || route[0] == "/") && rootPath != "" {
		http.HandleFunc(rootPath+"/", func(w http.ResponseWriter, r *http.Request) {
			var tempController PT = new(T)
			tempController.init(w, r)
			runRequest(w, r, tempController)
		})

		http.HandleFunc(rootPath, func(w http.ResponseWriter, r *http.Request) {
			var tempController PT = new(T)
			tempController.init(w, r)
			runRequest(w, r, tempController)
		})
		return
	}

	http.HandleFunc(rootPath+"/"+strings.Join(route, "/")+"/", func(w http.ResponseWriter, r *http.Request) {
		var tempController PT = new(T)
		tempController.init(w, r)
		runRequest(w, r, tempController)
	})

	http.HandleFunc(rootPath+"/"+strings.Join(route, "/"), func(w http.ResponseWriter, r *http.Request) {
		var tempController PT = new(T)
		tempController.init(w, r)
		runRequest(w, r, tempController)
	})
}
