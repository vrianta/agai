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
		sessions: make(map[string]Session.Struct),
		routes:   _routes,
	}
}

func (router *Struct) Handler(w http.ResponseWriter, r *http.Request) {

	sessionID := Session.GetSessionID(r)
	if sessionID == nil { // means user does not have any session with the server so creating a new clean guest session with the server
		Session := Session.New(w, r)
		sessionID = Session.StartSession()
		if sessionID != nil { // Successfuly started a New session without any error
			router.sessions[(*sessionID)] = *Session
			if _controller, ok := router.routes[r.URL.Path]; ok {
				Session.UpdateSession(&w, r)
				Session.ParseRequest()
				response := _controller.CallMethod(Session)
				if err := Session.RenderEngine.RenderTemplate(_controller.View, response); err != nil {
					Log.WriteLog("Error rendering template: " + err.Error())
					panic(err) // Panic if there is an error rendering the template
				} else {
					Session.RenderEngine.StartRender()
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

		if __session, ok := router.sessions[(*sessionID)]; ok {
			if _controller, ok := router.routes[r.URL.Path]; ok {
				__session.UpdateSession(&w, r)
				__session.ParseRequest()
				response := _controller.CallMethod(&__session)
				if err := __session.RenderEngine.RenderTemplate(_controller.View, response); err != nil {
					Log.WriteLog("Error rendering template: " + err.Error())
					panic(err) // Panic if there is an error rendering the template
				} else {
					__session.RenderEngine.StartRender()
				}
			} else {
				http.Error(w, "404 Error : Route not found ", 404)
			}
		} else { // server is not holding the session any more so creating a new guest session for the user
			__session := Session.New(w, r)
			sessionID = __session.StartSession()
			if sessionID != nil {
				router.sessions[(*sessionID)] = *__session
				if _controller, ok := router.routes[r.URL.Path]; ok {
					__session.ParseRequest()
					response := _controller.CallMethod(__session)
					if err := __session.RenderEngine.RenderTemplate(_controller.View, response); err != nil {
						Log.WriteLog("Error rendering template: " + err.Error())
						panic(err) // Panic if there is an error rendering the template
					} else {
						__session.RenderEngine.StartRender()
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

/*
 * Handling the Requests coming for the CSS Files specially
 */
func (s *Struct) CSSHandlers(w http.ResponseWriter, r *http.Request) {

	_file_path := "." + r.URL.Path // path of the file

	_file_record, file_record_ok := fileInfo[_file_path]

	info, err := os.Stat(_file_path)
	if err != nil {
		Log.WriteLog(err.Error())
		return
	}

	w.Header().Set("Content-Type", "text/css; charset=utf-8")

	if file_record_ok && _file_record.LastModified.Compare(info.ModTime()) != 0 { // file has not been updated so alreay no not make it do extra work

		w.Write([]byte(fileInfo[_file_path].Data))
		return
	}

	_file_data := Utils.ReadFromFile(_file_path)

	fileInfo[_file_path] = FileInfo{
		Uri:          _file_path,
		LastModified: info.ModTime(),
		Data:         _file_data,
	}

	w.Write([]byte(fileInfo[_file_path].Data))

}

/*
 * Handling the Requests coming for the Js Files specially
 */
func (s *Struct) JsHandler(w http.ResponseWriter, r *http.Request) {
	_file_path := "." + r.URL.Path // path of the file

	_file_record, file_record_ok := fileInfo[_file_path]

	info, err := os.Stat(_file_path)
	if err != nil {
		Log.WriteLog(err.Error())
		return
	}

	w.Header().Set("Content-Type", "text/css; charset=utf-8")

	if file_record_ok && _file_record.LastModified.Compare(info.ModTime()) != 0 { // file has not been updated so alreay no not make it do extra work

		w.Write([]byte(fileInfo[_file_path].Data))
		return
	}

	_file_data := Utils.ReadFromFile(_file_path)

	fileInfo[_file_path] = FileInfo{
		Uri:          _file_path,
		LastModified: info.ModTime(),
		Data:         _file_data,
	}

	w.Write([]byte(fileInfo[_file_path].Data))

}

// RemoveSession removes a session from the session manager
func (r *Struct) RemoveSession(sessionID string) {
	defer delete(r.sessions, sessionID)
}
