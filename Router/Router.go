package Router

import (
	"net/http"
	"os"
	"sync"

	"github.com/vrianta/Server/Log"
	Session "github.com/vrianta/Server/Session"
	"github.com/vrianta/Server/Utils"
)

// Constructor for Router
func New(_routes Type) *Struct {
	return &Struct{
		sessions: sync.Map{}, // Initialize the session map
		routes:   _routes,
	}
}

// Handler processes incoming HTTP requests and manages user sessions.
// It checks if the user has an existing session and handles session creation or validation.
// Based on the session and route, it invokes the appropriate controller method.
// Parameters:
// - w: The HTTP response writer.
// - r: The HTTP request.
func (router *Struct) Handler(w http.ResponseWriter, r *http.Request) {

	sessionID := Session.GetSessionID(r)
	if sessionID == nil { // means user does not have any session with the server so creating a new clean guest session with the server
		Session := Session.New(w, r)
		sessionID = Session.StartSession()
		if sessionID != nil { // Successfuly started a New session without any error
			router.sessions.Store((*sessionID), *Session)
			if _controller, ok := router.routes[r.URL.Path]; ok {
				Session.UpdateSession(w, r)
				Session.ParseRequest()
				response := _controller.CallMethod(Session)
				if err := Session.RenderEngine.RenderTemplate(_controller.View, response); err != nil {
					Log.WriteLog("Error rendering template: " + err.Error())
					panic(err) // Panic if there is an error rendering the template
				}
			} else {
				// WriteConsolef("Route not found for URL: %s \n", r.URL.Path)ss
				http.Error(w, "404 Error : Route not found ", 404)
			}
		} else {
			http.Error(w, "Server Error * Failed to Create the Session for the user", 500)
		}
	} else { // User has a session ID to begin with
		// checking if the session is valid or not means it is checking if the server also has the session or not
		// if the session is valid then it will just update the session with the latest value

		if __session, ok := router.sessions.Load((*sessionID)); ok {
			if _controller, ok := router.routes[r.URL.Path]; ok {
				sessionPtr := __session.(Session.Struct)
				sessionPtr.UpdateSession(w, r)
				sessionPtr.ParseRequest()
				response := _controller.CallMethod(&sessionPtr)
				if err := sessionPtr.RenderEngine.RenderTemplate(_controller.View, response); err != nil {
					Log.WriteLog("Error rendering template: " + err.Error())
					panic(err) // Panic if there is an error rendering the template
				}
			} else {
				http.Error(w, "404 Error : Route not found ", 404)
			}
		} else { // server is not holding the session any more so creating a new guest session for the user
			__session := Session.New(w, r)
			sessionID = __session.StartSession()
			if sessionID != nil {
				router.sessions.Store((*sessionID), *__session)
				if _controller, ok := router.routes[r.URL.Path]; ok {
					__session.ParseRequest()
					response := _controller.CallMethod(__session)
					if err := __session.RenderEngine.RenderTemplate(_controller.View, response); err != nil {
						Log.WriteLog("Error rendering template: " + err.Error())
						panic(err) // Panic if there is an error rendering the template
					}
				} else {
					http.Error(w, "404 Error : Route not found ", 404)
				}
			} else {
				http.Error(w, "Server Error * Failed to Create the Session for the user", 500)
			}
		}
	}
}

// StaticFileHandler serves static files with caching support.
// It checks if the file exists in the cache and serves it directly if the cache is valid.
// Otherwise, it reads the file from disk, caches it, and serves it.
// Parameters:
// - contentType: The MIME type of the file being served.
// Returns:
// - http.HandlerFunc: A handler function for serving static files.
func (s *Struct) StaticFileHandler(contentType string) http.HandlerFunc {
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

// RemoveSession removes a session from the session manager.
// It ensures the session is deleted after use.
// Parameters:
// - sessionID: The ID of the session to be removed.
func (r *Struct) RemoveSession(sessionID string) {
	defer r.sessions.Delete(sessionID) // Ensure the session is deleted after use
}
