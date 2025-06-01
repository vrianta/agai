package Router

import (
	"net/http"
	"os"

	"github.com/vrianta/Server/Log"
	Session "github.com/vrianta/Server/Session"
	"github.com/vrianta/Server/Utils"
)

// Constructor for Router
func New(_routes Type) *Struct {
	return &Struct{
		sessions: make(map[string]*Session.Struct), // Use a regular map for sessions
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
	var sess *Session.Struct
	var ok bool

	if sessionID == nil {
		// No session, create a new one
		sess = Session.New(w, r)
		sessionID = sess.StartSession()
		if sessionID == nil {
			http.Error(w, "Server Error * Failed to Create the Session for the user", http.StatusInternalServerError)
			return
		}
		router.sessionMutex.Lock()
		router.sessions[*sessionID] = sess
		router.sessionMutex.Unlock()
	} else {
		router.sessionMutex.RLock()
		sess, ok = router.sessions[*sessionID]
		router.sessionMutex.RUnlock()
		if !ok {
			// Session not found, create a new one
			sess = Session.New(w, r)
			sessionID = sess.StartSession()
			if sessionID == nil {
				http.Error(w, "Server Error * Failed to Create the Session for the user", http.StatusInternalServerError)
				return
			}
			router.sessionMutex.Lock()
			router.sessions[*sessionID] = sess
			router.sessionMutex.Unlock()
		}
	}

	// At this point, sess is valid
	if _controller, found := router.routes[r.URL.Path]; found {
		sess.UpdateSession(w, r)
		sess.ParseRequest()
		response := _controller.CallMethod(sess)
		if err := sess.RenderEngine.RenderTemplate(_controller.View, response); err != nil {
			Log.WriteLog("Error rendering template: " + err.Error())
			panic(err)
		}
	} else {
		http.Error(w, "404 Error : Route not found ", http.StatusNotFound)
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
	r.sessionMutex.Lock()
	defer r.sessionMutex.Unlock()
	delete(r.sessions, sessionID)
}

// Get Function to return all the Routes
func (r *Struct) Get() *Type {
	return &r.routes
}
