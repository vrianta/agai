package Router

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/vrianta/Server/Controller"
	"github.com/vrianta/Server/Log"
	Session "github.com/vrianta/Server/Session"
	"github.com/vrianta/Server/Utils"
)

// Constructor for Router
func InitRoutes(_routes *Routes) error {

	routes = *_routes

	return nil
}

// Handler processes incoming HTTP requests and manages user sessions.
// It checks if the user has an existing session and handles session creation or validation.
// Based on the session and route, it invokes the appropriate controller method.
// Parameters:
// - w: The HTTP response writer.
// - r: The HTTP request.
func Handler(w http.ResponseWriter, r *http.Request) {
	start := time.Now() // Start time measurement
	sessionID := Session.GetSessionID(r)
	var sess *Session.Struct
	var ok bool

	var tempController *Controller.Struct
	if _controller, found := routes[r.URL.Path]; found {
		tempController = _controller.Copy()
	} else {
		http.Error(w, "404 Error : Route not found ", http.StatusNotFound)
		return
	}

	if sessionID == nil {
		// No session, create a new one
		sess = Session.New(w, r)
		sessionID, err := Utils.GenerateSessionID()
		if err != nil {
			Log.WriteLog("Error generating session ID: " + err.Error())
			return
		}

		if sess.StartSession(&sessionID) == nil {
			http.Error(w, "Server Error * Failed to Create the Session for the user", http.StatusInternalServerError)
			return
		}
		Session.Store(&sessionID, sess)
	} else {
		sess, ok = Session.Get(sessionID)
		if !ok {
			// Session not found, create a new one
			sess = Session.New(w, r)
			sessionID, err := Utils.GenerateSessionID()
			if err != nil {
				Log.WriteLog("Error generating session ID: " + err.Error())
				return
			}

			if sess.StartSession(&sessionID) == nil {
				http.Error(w, "Server Error * Failed to Create the Session for the user", http.StatusInternalServerError)
				return
			}
			Session.Store(&sessionID, sess)
		}
	}

	sess.UpdateSession(w, r)
	sess.ParseRequest()
	response := tempController.CallMethod(sess)
	if err := tempController.Execute(response); err != nil {
		Log.WriteLog("Error rendering template: " + err.Error())
		panic(err)
	}

	duration := time.Since(start)
	log.Printf("Handler took %s to complete\n", duration)
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
		val, _ := fileInfo.Load(_filePath)
		cached, ok := val.(FileInfo)

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
		newRecord := FileInfo{
			Uri:          _filePath,
			LastModified: info.ModTime(),
			Data:         _fileData,
		}
		fileInfo.Store(_filePath, newRecord)
		w.Write([]byte(_fileData))
	}
}

// Get Function to return all the Routes
func Get() *Routes {
	return &routes
}

// return a list of all the views from the controllers
// loop throgh all the controllers and make a array of strings
func GetViews() []string {
	routerSize := len(routes)
	if routerSize < 1 {
		return nil
	}
	response := make([]string, routerSize)
	for _, controller := range routes {
		response = append(response, controller.View)
	}
	return response
}
