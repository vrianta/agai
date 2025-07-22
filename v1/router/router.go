package router

import (
	"net/http"
	"os"

	Controller "github.com/vrianta/agai/v1/controller"
	Log "github.com/vrianta/agai/v1/log"
	Utils "github.com/vrianta/agai/v1/utils"
)

/*
 * Create a New Router Object with Default route group example / is the default route for this or /api or /v1 etc
 */
func New(route string) *Struct {
	return &Struct{
		defaultRoute: route,
	}
}

/*
 * initialise Requests and Register the paths
 * Syntax - Router.New("").RegisterRoutes(
 *	NewRoute(path, controllerOnject),
 *  NewRoute(path, controllerOnject),
 * )
 * Example - Router.New("").RegisterRoutes(
 *	NewRoute("/home", homeObj),
 *  NewRoute("/list", listObj),
 * )
 */
func (_r *Struct) RegisterRoutes(_routes ...route) error {
	for _, rt := range _routes {
		routeTable[_r.defaultRoute+rt.path] = &rt.controllerObject
	}
	return nil
}

/*
 * Create Route Object
 */
func Route(path string, obj Controller.Context) route {
	return route{
		path:             path,
		controllerObject: obj,
	}
}

// Handler processes incoming HTTP requests and manages user sessions.
// It checks if the user has an existing session and handles session creation or validation.
// Based on the session and route, it invokes the appropriate controller method.
// Parameters:
// - w: The HTTP response writer.
// - r: The HTTP request.
func Handler(w http.ResponseWriter, r *http.Request) {

	var tempController *Controller.Context
	if _controller, found := routeTable[r.URL.Path]; found {
		tempController = _controller.Copy()
		tempController.Init(w, r)
	} else {
		http.Error(w, "404 Error : Route not found ", http.StatusNotFound)
		return
	}
}

// A Function to Create and Return

// StaticFileHandler serves static files with caching support.
// It checks if the file exists in the cache and serves it directly if the cache is valid.
// Otherwise, it reads the file from disk, caches it, and serves it.
// Parameters:
// - contentType: The MIME type of the file being served.
// Returns:
// - http.HandlerFunc: A handler function for serving static files.
func StaticFileHandler(contentType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_filePath := "." + r.URL.Path

		// Attempt to load from cache
		val, _ := fileCache.Load(_filePath)
		cached, ok := val.(FileCacheEntry)

		info, err := os.Stat(_filePath)
		if err != nil {
			Log.WriteLog(err.Error())
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", contentType)

		// If cached data exists and mod time matches then serve from cache
		if ok && cached.LastModified.Equal(info.ModTime()) {
			w.Write([]byte(cached.Data))
			return
		}

		// Read file from disk and cache it
		_fileData := Utils.ReadFromFile(_filePath)
		newRecord := FileCacheEntry{
			Uri:          _filePath,
			LastModified: info.ModTime(),
			Data:         _fileData,
		}
		fileCache.Store(_filePath, newRecord)
		w.Write([]byte(_fileData))
	}
}

// Get Function to return all the Routes
func GetRoutes() *routes {
	return &routeTable
}

// return a list of all the views from the controllers
// loop throgh all the controllers and make a array of strings
func ListViews() []string {
	routerSize := len(routeTable)
	if routerSize < 1 {
		return nil
	}
	response := make([]string, 0, routerSize)
	for _, controller := range routeTable {
		response = append(response, controller.View)
	}
	return response
}
